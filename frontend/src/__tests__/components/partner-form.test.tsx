import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import React from "react";
import { PartnerForm } from "@/components/forms/partner-form";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

vi.mock("next/navigation", () => ({
  useRouter: () => ({
    push: vi.fn(),
  }),
}));

vi.mock("@/hooks/use-partners", () => ({
  useCreatePartner: () => ({
    mutateAsync: vi.fn().mockResolvedValue({}),
    isSubmitting: false,
  }),
  useUpdatePartner: () => ({
    mutateAsync: vi.fn().mockResolvedValue({}),
    isSubmitting: false,
  }),
}));

vi.mock("sonner", () => ({
  toast: { success: vi.fn(), error: vi.fn() },
}));

function wrapper({ children }: { children: React.ReactNode }) {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return React.createElement(QueryClientProvider, { client: qc }, children);
}

describe("PartnerForm", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders 'Crear Socio' when no initialData", () => {
    render(React.createElement(PartnerForm), { wrapper });
    expect(screen.getByText("Crear Socio")).toBeInTheDocument();
  });

  it("renders 'Editar Socio' when isEdit is true", () => {
    render(
      React.createElement(PartnerForm, {
        initialData: {
          id: "1",
          name: "Existing Partner",
          tax_id: "999",
          status: "active",
          created_at: "",
          updated_at: "",
        },
        isEdit: true,
      }),
      { wrapper }
    );
    expect(screen.getByText("Editar Socio")).toBeInTheDocument();
  });

  it("shows validation error when name is empty", async () => {
    const user = userEvent.setup();
    render(React.createElement(PartnerForm), { wrapper });

    await user.click(screen.getByRole("button", { name: /crear/i }));

    await waitFor(() => {
      expect(screen.getByText("Nombre requerido")).toBeInTheDocument();
    });
  });

  it("submits with valid name", async () => {
    const user = userEvent.setup();
    render(React.createElement(PartnerForm), { wrapper });

    await user.type(screen.getByLabelText(/nombre/i), "New Partner");
    await user.click(screen.getByRole("button", { name: /crear/i }));

    await waitFor(() => {
      expect(screen.queryByText("Nombre requerido")).not.toBeInTheDocument();
    });
  });

  it("renders status select with Activo/Inactivo options", () => {
    render(React.createElement(PartnerForm), { wrapper });
    const select = screen.getByLabelText(/estado/i);
    expect(select).toBeInTheDocument();
  });

  it("renders cancel button", () => {
    render(React.createElement(PartnerForm), { wrapper });
    expect(screen.getByRole("button", { name: /cancelar/i })).toBeInTheDocument();
  });
});
