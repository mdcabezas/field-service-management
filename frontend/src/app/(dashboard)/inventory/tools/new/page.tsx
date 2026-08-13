"use client";

import { useRouter } from "next/navigation";
import { useCreateTool } from "@/hooks/use-tools";
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
  code: z.string().min(1, "Código es requerido"),
  status: z.string().min(1, "Estado es requerido"),
});

type FormData = z.infer<typeof schema>;

export default function NewToolPage() {
  const router = useRouter();
  const createTool = useCreateTool();

  const { register, handleSubmit, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: { name: "", code: "", status: "" },
  });

  const onSubmit = async (data: FormData) => {
    try {
      await createTool.mutateAsync(data);
      toast.success("Herramienta creada");
      router.back();
    } catch {
      toast.error("Error al guardar");
    }
  };

  return (
    <div className="space-y-4">
      <h1 className="font-mono text-2xl font-bold uppercase tracking-wider">
        Crear Herramienta
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
              <Label htmlFor="code">Código</Label>
              <Input id="code" {...register("code")} />
              {errors.code && <p className="text-sm text-red-500">{errors.code.message}</p>}
            </div>
            <div className="space-y-2">
              <Label htmlFor="status">Estado</Label>
              <select
                id="status"
                {...register("status")}
                className="w-full h-10 px-3 font-mono text-sm border-2 border-black bg-white rounded-none"
              >
                <option value="">Seleccionar...</option>
                <option value="available">Disponible</option>
                <option value="in_use">En Uso</option>
                <option value="maintenance">Mantenimiento</option>
                <option value="retired">Retirado</option>
              </select>
              {errors.status && <p className="text-sm text-red-500">{errors.status.message}</p>}
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
