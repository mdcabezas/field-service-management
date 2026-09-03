import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import React from "react";
import { TechRoleForm } from "@/components/forms/tech-role-form";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
}));
vi.mock("@/hooks/use-tech-roles", () => ({
  useCreateTechRole: () => ({ mutateAsync: vi.fn().mockResolvedValue({}), isSubmitting: false }),
  useUpdateTechRole: () => ({ mutateAsync: vi.fn().mockResolvedValue({}), isSubmitting: false }),
}));
vi.mock("sonner", () => ({ toast: { success: vi.fn(), error: vi.fn() } }));

function wrapper({ children }: { children: React.ReactNode }) {
  return React.createElement(QueryClientProvider, { client: new QueryClient({ defaultOptions: { queries: { retry: false } } }) }, children);
}

describe("TechRoleForm", () => {
  beforeEach(() => vi.clearAllMocks());

  it("renders 'Crear Rol' when no initialData", () => {
    render(React.createElement(TechRoleForm), { wrapper });
    expect(screen.getByText("Crear Rol")).toBeInTheDocument();
  });

  it("renders 'Editar Rol' in edit mode", () => {
    render(React.createElement(TechRoleForm, { initialData: { id: "1", code: "TECH", name: "Technician", active: true, created_at: "" }, isEdit: true }), { wrapper });
    expect(screen.getByText("Editar Rol")).toBeInTheDocument();
  });

  it("shows code validation error when empty", async () => {
    const user = userEvent.setup();
    render(React.createElement(TechRoleForm), { wrapper });
    await user.click(screen.getByRole("button", { name: /crear/i }));
    await waitFor(() => {
      expect(screen.getByText("Código requerido")).toBeInTheDocument();
    });
  });

  it("shows name validation error when empty", async () => {
    const user = userEvent.setup();
    render(React.createElement(TechRoleForm), { wrapper });
    await user.type(screen.getByLabelText(/código/i), "CODE");
    await user.click(screen.getByRole("button", { name: /crear/i }));
    await waitFor(() => {
      expect(screen.getByText("Nombre requerido")).toBeInTheDocument();
    });
  });

  it("renders active checkbox", () => {
    render(React.createElement(TechRoleForm), { wrapper });
    expect(screen.getByLabelText(/activo/i)).toBeInTheDocument();
  });

  it("renders cancel button", () => {
    render(React.createElement(TechRoleForm), { wrapper });
    expect(screen.getByRole("button", { name: /cancelar/i })).toBeInTheDocument();
  });
});
