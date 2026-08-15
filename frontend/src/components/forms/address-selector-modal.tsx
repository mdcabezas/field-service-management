"use client";

import { useState, useEffect, useCallback } from "react";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { AddressMap } from "@/components/map/address-map";
import { useAddresses, useCreateAddress } from "@/hooks/use-addresses";
import { useCreateCustomerAddress } from "@/hooks/use-customer-addresses";
import type { Address, CustomerAddressType } from "@/types/api";

interface AddressSelectorModalProps {
  customerId: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: (data: { address_id: string; name?: string; type: CustomerAddressType }) => void;
  initialAddressId?: string;
}

export function AddressSelectorModal({
  customerId,
  open,
  onOpenChange,
  onSuccess,
  initialAddressId,
}: AddressSelectorModalProps) {
  const [activeTab, setActiveTab] = useState<"select" | "create">("select");
  const [selectedAddress, setSelectedAddress] = useState<string | null>(initialAddressId || null);
  const [formData, setFormData] = useState({
    street: "",
    number: "",
    apartment: "",
    neighborhood: "",
    city: "",
    region: "",
    location_references: "",
    postal_code: "",
  });
  const [addressType, setAddressType] = useState<CustomerAddressType>("residential");
  const [name, setName] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const { data: addresses = [], isLoading: loadingAddresses } = useAddresses({ limit: 100 });
  const createAddress = useCreateAddress();
  const createCustomerAddress = useCreateCustomerAddress(customerId);

  useEffect(() => {
    if (initialAddressId) {
      setSelectedAddress(initialAddressId);
      setActiveTab("select");
    }
  }, [initialAddressId]);

  const handleCreateAddress = async () => {
    if (!formData.street || !formData.city) {
      setError("Calle y ciudad son requeridos");
      return;
    }

    setIsSubmitting(true);
    setError(null);

    try {
      const newAddress = await createAddress.mutateAsync({
        street: formData.street,
        number: formData.number || "",
        apartment: formData.apartment || "",
        neighborhood: formData.neighborhood || "",
        city: formData.city,
        region: formData.region || "",
        location_references: formData.location_references || "",
        postal_code: formData.postal_code || "",
      });

      await createCustomerAddress.mutateAsync({
        address_id: newAddress.id,
        type: addressType,
        name: name || undefined,
      });

      onSuccess({
        address_id: newAddress.id,
        name: name || undefined,
        type: addressType,
      });

      onOpenChange(false);
    } catch (err) {
      setError("Error al crear la dirección y vincularla al cliente");
      console.error(err);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleSelectExisting = async () => {
    if (!selectedAddress) {
      setError("Selecciona una dirección");
      return;
    }

    setIsSubmitting(true);
    setError(null);

    try {
      await createCustomerAddress.mutateAsync({
        address_id: selectedAddress,
        type: addressType,
        name: name || undefined,
      });

      onSuccess({
        address_id: selectedAddress,
        name: name || undefined,
        type: addressType,
      });

      onOpenChange(false);
    } catch (err) {
      setError("Error al vincular la dirección al cliente");
      console.error(err);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-4xl max-h-[90vh] overflow-hidden">
        <DialogHeader>
          <DialogTitle>Agregar Dirección al Cliente</DialogTitle>
        </DialogHeader>

        <DialogContent className="flex flex-col h-[70vh] p-0">
          <Tabs value={activeTab} onValueChange={setActiveTab} className="flex flex-col h-full">
            <TabsList className="border-b border-black bg-white">
              <TabsTrigger value="select">Seleccionar Existente</TabsTrigger>
              <TabsTrigger value="create">Crear Nueva</TabsTrigger>
            </TabsList>

            <TabsContent value="select" className="flex-1 overflow-hidden flex flex-col p-4">
              {loadingAddresses ? (
                <div className="flex items-center justify-center h-full font-mono text-gray-500">
                  Cargando direcciones...
                </div>
              ) : addresses.length === 0 ? (
                <div className="flex items-center justify-center h-full font-mono text-gray-500">
                  No hay direcciones disponibles. Crea una nueva en la pestaña "Crear Nueva".
                </div>
              ) : (
                <>
                  <div className="mb-4 font-mono text-sm">
                    <label className="block mb-2">Buscar dirección:</label>
                    <Input
                      placeholder="Filtrar por calle, ciudad..."
                      onChange={(e) => {
                        // Client-side filtering would go here
                      }}
                    />
                  </div>
                  <div className="flex-1 overflow-auto border-2 border-black rounded-lg">
                    <table className="w-full font-mono text-sm">
                      <thead className="bg-black text-white sticky top-0">
                        <tr>
                          <th className="p-2 text-left w-10">Sel.</th>
                          <th className="p-2 text-left">Dirección</th>
                          <th className="p-2 text-left">Ciudad</th>
                          <th className="p-2 text-left">Región</th>
                        </tr>
                      </thead>
                      <tbody>
                        {addresses.map((addr) => (
                          <tr
                            key={addr.id}
                            className={
                              `border-b border-gray-200 hover:bg-gray-50 transition-colors ${selectedAddress === addr.id ? "bg-gray-100" : ""}`
                            }
                          >
                            <td className="p-2">
                              <input
                                type="radio"
                                checked={selectedAddress === addr.id}
                                onChange={() => setSelectedAddress(addr.id)}
                                className="w-4 h-4 accent-black"
                              />
                            </td>
                            <td className="p-2">
                              {addr.street}
                              {addr.number && ` ${addr.number}`}
                              {addr.apartment && `, Depto ${addr.apartment}`}
                              {addr.neighborhood && `, ${addr.neighborhood}`}
                              {addr.location_references && ` (${addr.location_references})`}
                            </td>
                            <td className="p-2">{addr.city}</td>
                            <td className="p-2">{addr.region}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </>
              )}

              {error && (
                <div className="mt-4 p-3 bg-red-50 border-2 border-red-500 text-red-700 font-mono text-sm">
                  {error}
                </div>
              )}
            </TabsContent>

            <TabsContent value="create" className="flex-1 overflow-auto p-4 space-y-4">
              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="street">Calle *</Label>
                  <Input
                    id="street"
                    value={formData.street}
                    onChange={(e) => setFormData({ ...formData, street: e.target.value })}
                    placeholder="Av. Providencia"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="number">Número</Label>
                  <Input
                    id="number"
                    value={formData.number}
                    onChange={(e) => setFormData({ ...formData, number: e.target.value })}
                    placeholder="123"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="apartment">Depto/Oficina</Label>
                  <Input
                    id="apartment"
                    value={formData.apartment}
                    onChange={(e) => setFormData({ ...formData, apartment: e.target.value })}
                    placeholder="Depto 501"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="neighborhood">Barrio/Sector</Label>
                  <Input
                    id="neighborhood"
                    value={formData.neighborhood}
                    onChange={(e) => setFormData({ ...formData, neighborhood: e.target.value })}
                    placeholder="Providencia"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="city">Ciudad *</Label>
                  <Input
                    id="city"
                    value={formData.city}
                    onChange={(e) => setFormData({ ...formData, city: e.target.value })}
                    placeholder="Santiago"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="region">Región</Label>
                  <Input
                    id="region"
                    value={formData.region}
                    onChange={(e) => setFormData({ ...formData, region: e.target.value })}
                    placeholder="RM"
                  />
                </div>
                <div className="space-y-2 md:col-span-2">
                  <Label htmlFor="neighborhood">Referencia de Ubicación</Label>
                  <Input
                    id="location_references"
                    value={formData.location_references}
                    onChange={(e) => setFormData({ ...formData, location_references: e.target.value })}
                    placeholder="Cerca de metro Los Leones"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="postal_code">Código Postal</Label>
                  <Input
                    id="postal_code"
                    value={formData.postal_code}
                    onChange={(e) => setFormData({ ...formData, postal_code: e.target.value })}
                    placeholder="7500000"
                  />
                </div>
              </div>

              <div className="space-y-2 pt-4 border-t border-black">
                <div className="grid gap-2 md:grid-cols-2">
                  <div className="space-y-2">
                    <Label htmlFor="type">Tipo de Dirección *</Label>
                    <Select value={addressType} onValueChange={(v) => setAddressType(v as CustomerAddressType)}>
                      <SelectTrigger id="type">
                        <SelectValue placeholder="Seleccionar tipo" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="residential">Residencial</SelectItem>
                        <SelectItem value="commercial">Comercial</SelectItem>
                        <SelectItem value="industrial">Industrial</SelectItem>
                        <SelectItem value="institutional">Institucional</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="name">Nombre (opcional)</Label>
                    <Input
                      id="name"
                      value={name}
                      onChange={(e) => setName(e.target.value)}
                      placeholder="Casa matriz, Sucursal centro, etc."
                    />
                  </div>
                </div>
              </div>

              {error && (
                <div className="p-3 bg-red-50 border-2 border-red-500 text-red-700 font-mono text-sm">
                  {error}
                </div>
              )}
            </TabsContent>
          </Tabs>

          <DialogFooter className="border-t border-black bg-white p-4">
            <div className="w-full flex justify-end gap-2">
              <Button variant="outline" onClick={() => onOpenChange(false)} disabled={isSubmitting}>
                Cancelar
              </Button>
              <Button
                onClick={activeTab === "select" ? handleSelectExisting : handleCreateAddress}
                disabled={isSubmitting || (activeTab === "select" && !selectedAddress)}
              >
                {isSubmitting ? "Guardando..." : activeTab === "select" ? "Usar esta dirección" : "Crear y agregar"}
              </Button>
            </div>
          </DialogFooter>
        </DialogContent>
      </DialogContent>
    </Dialog>
  );
}