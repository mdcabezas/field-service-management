"use client";

import { useState, useEffect } from "react";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { AddressSelectorModal } from "@/components/forms/address-selector-modal";
import { useCustomerAddresses, useCreateCustomerAddress, useUpdateCustomerAddress, useDeleteCustomerAddress } from "@/hooks/use-customer-addresses";
import type { CustomerAddress, CustomerAddressType } from "@/types/api";

interface CustomerAddressFormProps {
  customerId: string;
  initialData?: CustomerAddress;
  isEdit: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

export function CustomerAddressForm({
  customerId,
  initialData,
  isEdit,
  onClose,
  onSuccess,
}: CustomerAddressFormProps) {
  const [showAddressSelector, setShowAddressSelector] = useState(false);
  const [selectedAddressId, setSelectedAddressId] = useState<string | null>(initialData?.address_id || null);
  const [type, setType] = useState<CustomerAddressType>(initialData?.type as CustomerAddressType || "residential");
  const [name, setName] = useState(initialData?.name || "");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const createCustomerAddress = useCreateCustomerAddress(customerId);
  const updateCustomerAddress = useUpdateCustomerAddress(customerId);
  const deleteCustomerAddress = useDeleteCustomerAddress(customerId);

  // When initialData changes (e.g., when editing), update form state
  useEffect(() => {
    if (initialData) {
      setSelectedAddressId(initialData.address_id || null);
      setType((initialData.type as CustomerAddressType) || "residential");
      setName(initialData.name || "");
    }
  }, [initialData]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (!selectedAddressId) {
      setError("Debes seleccionar una dirección");
      return;
    }

    setIsSubmitting(true);

    try {
      if (isEdit && initialData) {
        await updateCustomerAddress.mutateAsync({
          id: initialData.id,
          data: {
            address_id: selectedAddressId,
            type,
            name: name || undefined,
          },
        });
      } else {
        await createCustomerAddress.mutateAsync({
          address_id: selectedAddressId,
          type,
          name: name || undefined,
        });
      }

      onSuccess();
      onClose();
    } catch (err) {
      setError("Error al guardar la dirección");
      console.error(err);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDelete = async () => {
    if (!isEdit || !initialData) return;

    if (!confirm("¿Estás seguro de eliminar esta dirección del cliente?")) return;

    try {
      await deleteCustomerAddress.mutateAsync(initialData.id);
      onSuccess();
      onClose();
    } catch (err) {
      setError("Error al eliminar la dirección");
      console.error(err);
    }
  };

  return (
    <Dialog open={true} onOpenChange={onClose}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>{isEdit ? "Editar Dirección" : "Agregar Dirección"}</DialogTitle>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4 p-4">
          <div className="space-y-2">
            <Label htmlFor="address">Dirección *</Label>
            <Button
              type="button"
              variant="outline"
              className="w-full justify-start"
              onClick={() => setShowAddressSelector(true)}
            >
              {selectedAddressId ? "Cambiar dirección" : "Seleccionar dirección"}
            </Button>
            {error && (
              <p className="text-sm text-red-500">{error}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="type">Tipo *</Label>
            <Select value={type} onValueChange={(v) => setType(v as CustomerAddressType)}>
              <SelectTrigger>
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

          {error && (
            <div className="p-3 bg-red-50 border-2 border-red-500 text-red-700 font-mono text-sm">
              {error}
            </div>
          )}

          <DialogFooter className="flex justify-end gap-2">
            <Button type="button" variant="outline" onClick={onClose} disabled={isSubmitting}>
              Cancelar
            </Button>
            {isEdit && (
              <Button type="button" variant="destructive" onClick={handleDelete} disabled={isSubmitting}>
                Eliminar
              </Button>
            )}
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? "Guardando..." : isEdit ? "Actualizar" : "Agregar"}
            </Button>
          </DialogFooter>
        </form>

        <AddressSelectorModal
          customerId={customerId}
          open={showAddressSelector}
          onOpenChange={setShowAddressSelector}
          onSuccess={(data) => {
            setSelectedAddressId(data.address_id);
            setType(data.type);
            if (data.name) setName(data.name);
            setShowAddressSelector(false);
          }}
          initialAddressId={selectedAddressId || undefined}
        />
      </DialogContent>
    </Dialog>
  );
}