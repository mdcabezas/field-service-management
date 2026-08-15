"use client";

import { useState, useMemo } from "react";
import { DataTable } from "@/components/data-table/data-table";
import { type ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { useCustomerAddresses, useCreateCustomerAddress, useDeleteCustomerAddress } from "@/hooks/use-customer-addresses";
import { useAddresses } from "@/hooks/use-addresses";
import type { CustomerAddress, CustomerAddressType, Address } from "@/types/api";

interface CustomerAddressesTableProps {
  customerId: string;
}

const typeLabels: Record<string, string> = {
  residential: "Residencial",
  commercial: "Comercial",
  industrial: "Industrial",
  institutional: "Institucional",
};

export function CustomerAddressesTable({ customerId }: CustomerAddressesTableProps) {
  const { data: customerAddresses = [], isLoading, refetch } = useCustomerAddresses(customerId);
  const { data: geocodingAddresses = [] } = useAddresses({ limit: 200 });
  const deleteCustomerAddress = useDeleteCustomerAddress(customerId);
  const createCustomerAddress = useCreateCustomerAddress(customerId);

  const [showDialog, setShowDialog] = useState(false);
  const [selectedAddressId, setSelectedAddressId] = useState<string>("");
  const [addressType, setAddressType] = useState<CustomerAddressType>("residential");
  const [name, setName] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState("");

  const addressMap = useMemo(() => {
    const map = new Map<string, Address>();
    geocodingAddresses.forEach((a) => map.set(a.id, a));
    return map;
  }, [geocodingAddresses]);

  const filteredAddresses = useMemo(() => {
    if (!searchQuery) return geocodingAddresses;
    const q = searchQuery.toLowerCase();
    return geocodingAddresses.filter(
      (a) =>
        a.street.toLowerCase().includes(q) ||
        a.city.toLowerCase().includes(q) ||
        (a.number && a.number.toLowerCase().includes(q))
    );
  }, [geocodingAddresses, searchQuery]);

  const handleOpenDialog = () => {
    setSelectedAddressId("");
    setAddressType("residential");
    setName("");
    setError(null);
    setSearchQuery("");
    setShowDialog(true);
  };

  const handleDelete = async (id: string) => {
    if (!confirm("¿Estás seguro de eliminar esta dirección del cliente?")) return;
    try {
      await deleteCustomerAddress.mutateAsync(id);
      refetch();
    } catch (err) {
      console.error("Error deleting address:", err);
    }
  };

  const handleSubmit = async () => {
    if (!selectedAddressId) {
      setError("Selecciona una dirección");
      return;
    }
    setIsSubmitting(true);
    setError(null);
    try {
      await createCustomerAddress.mutateAsync({
        address_id: selectedAddressId,
        type: addressType,
        name: name || undefined,
      });
      setShowDialog(false);
      refetch();
    } catch (err) {
      setError("Error al vincular la dirección al cliente");
      console.error(err);
    } finally {
      setIsSubmitting(false);
    }
  };

  const columns: ColumnDef<CustomerAddress, unknown>[] = [
    {
      accessorKey: "address",
      header: "Dirección",
      cell: ({ row }) => {
        const ca = row.original;
        const addr = addressMap.get(ca.address_id);
        return (
          <div className="font-mono text-sm">
            <div className="font-medium">
              {addr ? `${addr.street}${addr.number ? ` ${addr.number}` : ""}` : ca.address_id?.slice(0, 8) + "..."}
            </div>
            {addr && (
              <div className="text-gray-500">
                {addr.neighborhood && `${addr.neighborhood}, `}{addr.city}{addr.region ? `, ${addr.region}` : ""}
              </div>
            )}
          </div>
        );
      },
    },
    {
      accessorKey: "type",
      header: "Tipo",
      cell: ({ row }) => (
        <Badge variant="outline" className="font-mono text-xs">
          {typeLabels[row.original.type as CustomerAddressType] || row.original.type}
        </Badge>
      ),
    },
    {
      accessorKey: "name",
      header: "Nombre",
      cell: ({ row }) => (
        <span className="font-mono text-sm">{row.original.name || "—"}</span>
      ),
    },
    {
      id: "actions",
      header: "Acciones",
      cell: ({ row }) => (
        <Button
          variant="destructive"
          size="sm"
          onClick={() => handleDelete(row.original.id)}
          className="h-8 px-3 font-mono text-xs"
        >
          Eliminar
        </Button>
      ),
    },
  ];

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64 font-mono text-gray-500">
        Cargando direcciones...
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="font-mono text-lg font-semibold">Direcciones del Cliente</h3>
        <Button onClick={handleOpenDialog}>+ Agregar Dirección</Button>
      </div>

      <DataTable
        columns={columns}
        data={customerAddresses}
        searchPlaceholder="Buscar por nombre..."
        searchColumn="name"
      />

      <Dialog open={showDialog} onOpenChange={setShowDialog}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>Agregar Dirección</DialogTitle>
          </DialogHeader>

          <div className="space-y-4 p-4">
            <div className="space-y-2">
              <Label>Dirección *</Label>
              <Input
                placeholder="Buscar por calle, ciudad..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
              <div className="border-2 border-black rounded-lg max-h-48 overflow-auto">
                {filteredAddresses.length === 0 ? (
                  <div className="p-4 text-center font-mono text-sm text-gray-500">
                    No hay direcciones disponibles
                  </div>
                ) : (
                  filteredAddresses.map((addr) => (
                    <button
                      key={addr.id}
                      type="button"
                      className={`w-full text-left p-2 font-mono text-sm border-b border-gray-200 hover:bg-gray-50 transition-colors ${selectedAddressId === addr.id ? "bg-gray-100 font-bold" : ""}`}
                      onClick={() => setSelectedAddressId(addr.id)}
                    >
                      {addr.street}{addr.number ? ` ${addr.number}` : ""}
                      <span className="text-gray-500 ml-2">— {addr.city}</span>
                    </button>
                  ))
                )}
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>Tipo *</Label>
                <Select value={addressType} onValueChange={(v) => setAddressType(v as CustomerAddressType)}>
                  <SelectTrigger>
                    <SelectValue />
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
                <Label>Nombre (opcional)</Label>
                <Input
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="Casa matriz, etc."
                />
              </div>
            </div>

            {error && (
              <div className="p-3 bg-red-50 border-2 border-red-500 text-red-700 font-mono text-sm">
                {error}
              </div>
            )}
          </div>

          <DialogFooter className="border-t border-black bg-white p-4">
            <Button variant="outline" onClick={() => setShowDialog(false)} disabled={isSubmitting}>
              Cancelar
            </Button>
            <Button onClick={handleSubmit} disabled={isSubmitting || !selectedAddressId}>
              {isSubmitting ? "Guardando..." : "Agregar"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
