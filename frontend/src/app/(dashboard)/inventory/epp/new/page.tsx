"use client";

import { useRouter } from "next/navigation";
import { useCreateEPPItem } from "@/hooks/use-epp-items";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

const schema = z.object({
  name: z.string().min(1, "Nombre es requerido"),
  type: z.string().min(1, "Tipo es requerido"),
  lifecycle: z.string().min(1, "Ciclo de vida es requerido"),
});

type FormData = z.infer<typeof schema>;

export default function NewEPPPage() {
  const router = useRouter();
  const createEPP = useCreateEPPItem();

  const { register, handleSubmit, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: { name: "", type: "", lifecycle: "" },
  });

  const onSubmit = async (data: FormData) => {
    try {
      await createEPP.mutateAsync(data);
      toast.success("EPP creado");
      router.back();
    } catch {
      toast.error("Error al guardar");
    }
  };

  return (
    <div className="space-y-4">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Crear EPP
      </h1>
      <Card className="max-w-md">
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">Nombre</Label>
              <Input id="name" {...register("name")} />
              {errors.name && <p className="text-sm text-red-500">{errors.name.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="type">Tipo</Label>
              <select
                id="type"
                {...register("type")}
                className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
              >
                <option value="">Seleccionar...</option>
                <option value="head">Cabeza</option>
                <option value="hands">Manos</option>
                <option value="body">Cuerpo</option>
                <option value="eyes">Ojos</option>
                <option value="ears">Oídos</option>
                <option value="feet">Pies</option>
                <option value="respiratory">Respiratorio</option>
              </select>
              {errors.type && <p className="text-sm text-red-500">{errors.type.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="lifecycle">Ciclo de Vida</Label>
              <select
                id="lifecycle"
                {...register("lifecycle")}
                className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
              >
                <option value="">Seleccionar...</option>
                <option value="disposable">Desechable</option>
                <option value="reusable">Reutilizable</option>
              </select>
              {errors.lifecycle && <p className="text-sm text-red-500">{errors.lifecycle.message}</p>}
            </div>
            <div className="flex gap-2">
              <Button type="submit">Crear</Button>
              <Button type="button" variant="outline" onClick={() => router.back()}>Cancelar</Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
