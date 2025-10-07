import { apiGet, apiPost, apiPut, apiDel, useApiClient } from "./_client";
import type { Page } from "./_base";

export interface Template {
  id: number;
  name: string;
  description: string;
  content: string;
  created_at?: string;
  updated_at?: string;
}

export function useTemplateApi() {
  const { baseURL } = useApiClient();

  const listTemplates = (page = 1, page_size = 20, q = "") =>
    apiGet<Page<Template>>("templates", {
      page,
      page_size,
      q: q || undefined,
    });

  const getTemplate = (id: number | string) =>
    apiGet<Template>(`templates/${id}`);

  const createTemplate = (data: Partial<Template>) =>
    apiPost<Template>("templates", data);

  const updateTemplate = (id: number | string, data: Partial<Template>) =>
    apiPut<Template>(`templates/${id}`, data);

  const deleteTemplate = (id: number | string) =>
    apiDel<void>(`templates/${id}`);

  return {
    baseURL,
    listTemplates,
    getTemplate,
    createTemplate,
    updateTemplate,
    deleteTemplate,
  };
}
