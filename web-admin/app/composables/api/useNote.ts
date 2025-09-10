import {
  apiGet,
  apiPost,
  apiPut,
  apiPatch,
  apiDel,
  useApiClient,
} from "./_client";
import type { Page } from "./_base";

export interface Note {
  id: string;
  title: string;
  content: string;
  category?: string;
  status?: "draft" | "published" | "archived";
  priority?: "low" | "medium" | "high";
  tags?: string[];
  created_at: string;
  updated_at: string;
  author_id?: string;
  author_name?: string;
}

export function useNoteApi() {
  const { baseURL } = useApiClient();

  const listNotes = (
    page = 1,
    page_size = 20,
    filters?: {
      category?: string;
      status?: string;
      priority?: string;
      search?: string;
      author_id?: string;
    }
  ) => {
    const query: Record<string, any> = { page, page_size };
    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== "") {
          query[key] = value;
        }
      });
    }
    return apiGet<Page<Note>>("notes", query);
  };

  const getNote = (id: string) => apiGet<Note>(`notes/${id}`);

  const createNote = (data: Partial<Note>) => apiPost<Note>("notes", data);

  const updateNote = (id: string, data: Partial<Note>) =>
    apiPut<Note>(`notes/${id}`, data);

  const deleteNote = (id: string) => apiDel<void>(`notes/${id}`);

  const searchNotes = (query: string, limit = 10) =>
    apiGet<Note[]>("notes/search", { q: query, limit });

  const getNotesByCategory = (category: string, page = 1, page_size = 20) =>
    apiGet<Page<Note>>("notes", { category, page, page_size });

  const getNotesByTag = (tag: string, page = 1, page_size = 20) =>
    apiGet<Page<Note>>("notes", { tag, page, page_size });

  const duplicateNote = (id: string) => apiPost<Note>(`notes/${id}/duplicate`);

  const archiveNote = (id: string) =>
    apiPatch<Note>(`notes/${id}`, { status: "archived" });

  const publishNote = (id: string) =>
    apiPatch<Note>(`notes/${id}`, { status: "published" });

  return {
    baseURL,
    listNotes,
    getNote,
    createNote,
    updateNote,
    deleteNote,
    searchNotes,
    getNotesByCategory,
    getNotesByTag,
    duplicateNote,
    archiveNote,
    publishNote,
  };
}
