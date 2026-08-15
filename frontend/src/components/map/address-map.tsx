"use client";

import { useEffect, useRef, useState } from "react";
import * as maplibregl from "maplibre-gl";
import "maplibre-gl/dist/maplibre-gl.css";
import type { Address } from "@/types/api";

interface AddressMapProps {
  addresses: Address[];
  onSelect: (address: Address) => void;
  onCreate: (lng: number, lat: number) => void;
  initialCenter?: [number, number];
  selectedAddressId?: string;
}

const DEFAULT_CENTER: [number, number] = [-70.65, -33.45];
const DEFAULT_ZOOM = 12;

export function AddressMap({
  addresses,
  onSelect,
  onCreate,
  initialCenter = DEFAULT_CENTER,
  selectedAddressId,
}: AddressMapProps) {
  const mapContainerRef = useRef<HTMLDivElement>(null);
  const mapRef = useRef<maplibregl.Map | null>(null);
  const markersRef = useRef<maplibregl.Marker[]>([]);
  const [mapLoaded, setMapLoaded] = useState(false);
  const [mapStyle, setMapStyle] = useState("dark");

  useEffect(() => {
    if (!mapContainerRef.current || mapRef.current) return;

    const map = new maplibregl.Map({
      container: mapContainerRef.current,
      style: `https://demotiles.maplibre.org/style.json`,
      center: initialCenter,
      zoom: DEFAULT_ZOOM,
      attributionControl: false,
    });

    map.addControl(new maplibregl.NavigationControl(), "top-right");
    map.addControl(new maplibregl.ScaleControl({ unit: "metric" }), "bottom-left");

    map.on("load", () => {
      setMapLoaded(true);
      addAddressMarkers(map);
    });

    map.on("click", (e) => {
      if (mapLoaded) {
        onCreate(e.lngLat.lng, e.lngLat.lat);
      }
    });

    mapRef.current = map;

    return () => {
      markersRef.current.forEach((marker) => marker.remove());
      markersRef.current = [];
      map.remove();
      mapRef.current = null;
      setMapLoaded(false);
    };
  }, [initialCenter, addresses, onSelect, onCreate, mapLoaded]);

  const addAddressMarkers = (map: maplibregl.Map) => {
    markersRef.current.forEach((marker) => marker.remove());
    markersRef.current = [];

    addresses.forEach((address) => {
      if (address.geom && Array.isArray(address.geom) && address.geom.length === 2) {
        const [lng, lat] = address.geom as [number, number];
        const isSelected = address.id === selectedAddressId;

        const el = document.createElement("div");
        el.className = `address-marker ${isSelected ? "selected" : ""}`;
        el.innerHTML = `
          <div class="marker-inner ${isSelected ? "selected" : ""}">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z" />
              <circle cx="12" cy="10" r="3" />
            </svg>
          </div>
        `;
        el.style.cssText = `
          cursor: pointer;
          transition: transform 0.2s;
          z-index: ${isSelected ? 1000 : 100};
        `;

        el.addEventListener("mouseenter", () => {
          el.style.transform = "scale(1.2)";
        });
        el.addEventListener("mouseleave", () => {
          el.style.transform = "scale(1)";
        });
        el.addEventListener("click", (e) => {
          e.stopPropagation();
          onSelect(address);
        });

        const marker = new maplibregl.Marker({ element: el })
          .setLngLat([lng, lat])
          .addTo(map);

        markersRef.current.push(marker);
      }
    });
  };

  const updateMarkers = () => {
    if (mapRef.current && mapLoaded) {
      addAddressMarkers(mapRef.current);
    }
  };

  useEffect(() => {
    updateMarkers();
  }, [addresses, selectedAddressId, mapLoaded]);

  return (
    <div
      ref={mapContainerRef}
      className="address-map"
      style={{
        width: "100%",
        height: "400px",
        borderRadius: "8px",
        overflow: "hidden",
        border: "2px solid black",
        backgroundColor: "#1a1a1a",
      }}
    >
      {!mapLoaded && (
        <div
          className="map-loading"
          style={{
            position: "absolute",
            inset: 0,
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            backgroundColor: "#1a1a1a",
            color: "#888",
            fontFamily: "monospace",
            fontSize: "14px",
            zIndex: 10,
          }}
        >
          Cargando mapa...
        </div>
      )}
    </div>
  );
}