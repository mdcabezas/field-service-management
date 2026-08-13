"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { authApi } from "@/lib/api";
import { useAuthStore } from "@/stores/auth-store";

const loginSchema = z.object({
  employeeNumber: z.string().min(1, "Número de empleado requerido"),
  password: z.string().min(1, "Contraseña requerida"),
});

type LoginFormData = z.infer<typeof loginSchema>;

export default function LoginPage() {
  const router = useRouter();
  const { login } = useAuthStore();
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
  });

  const onSubmit = async (data: LoginFormData) => {
    setError(null);
    setIsLoading(true);

    try {
      const tokens = await authApi.login(data.employeeNumber, data.password);

      // Store tokens in Zustand + cookie (cookie needed for middleware)
      authApi.setAuthTokens(tokens.access_token, tokens.refresh_token);

      // Get user info (now with Authorization header)
      const userInfo = await authApi.me();
      login(userInfo.employee_number, userInfo.role as "admin" | "manager" | "supervisor" | "technician");

      router.push("/");
    } catch (err) {
      setError("Credenciales inválidas");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-white">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>LOCALIS FSM</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="employeeNumber">Número de Empleado</Label>
              <Input
                id="employeeNumber"
                type="text"
                placeholder="Ej: 1001"
                {...register("employeeNumber")}
              />
              {errors.employeeNumber && (
                <p className="text-sm text-red-600">{errors.employeeNumber.message}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="password">Contraseña</Label>
              <Input
                id="password"
                type="password"
                placeholder="••••••••"
                {...register("password")}
              />
              {errors.password && (
                <p className="text-sm text-red-600">{errors.password.message}</p>
              )}
            </div>

            {error && (
              <p className="text-sm text-red-600 font-mono">{error}</p>
            )}

            <Button type="submit" className="w-full" disabled={isLoading}>
              {isLoading ? "Ingresando..." : "Iniciar Sesión"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
