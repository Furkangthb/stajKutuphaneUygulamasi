import { useEffect, useState } from "react";
import { api, type Book } from "../api";

const empty: Omit<Book, "id"> = {
  title: "",
  author: "",
  isbn: "",
  genre: "",
  published_year: undefined,
  available: true,
  description: "",
};

export default function ManageBooksPage() {
  const [books, setBooks] = useState<Book[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [modal, setModal] = useState<{ open: boolean; editing: Book | null }>({ open: false, editing: null });
  const [form, setForm] = useState<Omit<Book, "id">>(empty);
  const [saving, setSaving] = useState(false);
  const [deleting, setDeleting] = useState<number | null>(null);
  const [toast, setToast] = useState("");
  const [search, setSearch] = useState("");

  useEffect(() => {
    load();
  }, []);

  const load = () => {
    setLoading(true);
    api.getBooks()
      .then(setBooks)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  };

  const showToast = (msg: string) => {
    setToast(msg);
    setTimeout(() => setToast(""), 3000);
  };

  const openAdd = () => {
    setForm(empty);
    setModal({ open: true, editing: null });
  };

  const openEdit = (book: Book) => {
    setForm({ title: book.title, author: book.author, isbn: book.isbn ?? "", genre: book.genre ?? "", published_year: book.published_year, available: book.available, description: book.description ?? "" });
    setModal({ open: true, editing: book });
  };

  const save = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      if (modal.editing) {
        const updated = await api.updateBook(modal.editing.id, form);
        setBooks((bs) => bs.map((b) => (b.id === modal.editing!.id ? updated : b)));
        showToast("Kitap güncellendi");
      } else {
        const added = await api.addBook(form);
        setBooks((bs) => [added, ...bs]);
        showToast("Kitap eklendi");
      }
      setModal({ open: false, editing: null });
    } catch (e: unknown) {
      showToast(e instanceof Error ? e.message : "İşlem başarısız");
    } finally {
      setSaving(false);
    }
  };

  const del = async (book: Book) => {
    if (!confirm(`"${book.title}" kitabını silmek istediğinizden emin misiniz?`)) return;
    setDeleting(book.id);
    try {
      await api.deleteBook(book.id);
      setBooks((bs) => bs.filter((b) => b.id !== book.id));
      showToast("Kitap silindi");
    } catch (e: unknown) {
      showToast(e instanceof Error ? e.message : "Silme başarısız");
    } finally {
      setDeleting(null);
    }
  };

  const filtered = books.filter(
    (b) =>
      b.title.toLowerCase().includes(search.toLowerCase()) ||
      b.author.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="h-full flex flex-col">
      <div className="px-8 pt-8 pb-6 flex items-start justify-between gap-4">
        <div>
          <h1 className="font-display text-3xl mb-1" style={{ color: "var(--foreground)" }}>
            Kitap Yönetimi
          </h1>
          <p className="text-sm" style={{ color: "var(--muted-foreground)" }}>
            {books.length} kitap
          </p>
        </div>
        <button
          onClick={openAdd}
          className="flex items-center gap-2 px-4 py-2.5 rounded-xl text-sm font-semibold cursor-pointer transition-all"
          style={{ backgroundColor: "var(--primary)", color: "var(--primary-foreground)" }}
        >
          <span className="text-base leading-none">+</span> Kitap Ekle
        </button>
      </div>

      <div className="px-8 pb-4">
        <div className="relative max-w-sm">
          <svg className="absolute left-3 top-1/2 -translate-y-1/2 pointer-events-none" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="var(--muted-foreground)" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <circle cx="11" cy="11" r="8" /><line x1="21" y1="21" x2="16.65" y2="16.65" />
          </svg>
          <input
            type="text"
            placeholder="Kitap veya yazar ara…"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full pl-9 pr-4 py-2.5 rounded-xl text-sm outline-none"
            style={{ backgroundColor: "var(--card)", border: "1.5px solid var(--border)", color: "var(--foreground)" }}
          />
        </div>
      </div>

      <div className="flex-1 overflow-y-auto px-8 pb-8">
        {loading ? (
          <LoadingSkeleton />
        ) : error ? (
          <p className="text-sm" style={{ color: "#DC2626" }}>{error}</p>
        ) : (
          <div className="rounded-2xl overflow-hidden" style={{ backgroundColor: "var(--card)", border: "1px solid var(--border)" }}>
            <table className="w-full text-sm">
              <thead>
                <tr style={{ borderBottom: "1px solid var(--border)" }}>
                  {["Başlık", "Yazar", "Tür", "Yıl", "Durum", ""].map((h) => (
                    <th
                      key={h}
                      className="px-5 py-3 text-left font-semibold text-xs uppercase tracking-wider"
                      style={{ color: "var(--muted-foreground)" }}
                    >
                      {h}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {filtered.map((book, i) => (
                  <tr
                    key={book.id}
                    style={{ borderBottom: i < filtered.length - 1 ? "1px solid var(--border)" : "none" }}
                    className="hover:bg-[var(--muted)] transition-colors duration-100"
                  >
                    <td className="px-5 py-3 font-semibold" style={{ color: "var(--foreground)" }}>{book.title}</td>
                    <td className="px-5 py-3" style={{ color: "var(--muted-foreground)" }}>{book.author}</td>
                    <td className="px-5 py-3" style={{ color: "var(--muted-foreground)" }}>{book.genre ?? "—"}</td>
                    <td className="px-5 py-3" style={{ color: "var(--muted-foreground)" }}>{book.published_year ?? "—"}</td>
                    <td className="px-5 py-3">
                      <span
                        className="px-2 py-0.5 rounded-full text-xs font-bold"
                        style={{ backgroundColor: book.available ? "#D1FAE5" : "#FEE2E2", color: book.available ? "#065F46" : "#991B1B" }}
                      >
                        {book.available ? "Mevcut" : "Dolu"}
                      </span>
                    </td>
                    <td className="px-5 py-3">
                      <div className="flex items-center gap-2">
                        <button
                          onClick={() => openEdit(book)}
                          className="px-3 py-1.5 rounded-lg text-xs font-semibold cursor-pointer transition-all"
                          style={{ backgroundColor: "var(--secondary)", color: "var(--secondary-foreground)" }}
                        >
                          Düzenle
                        </button>
                        <button
                          onClick={() => del(book)}
                          disabled={deleting === book.id}
                          className="px-3 py-1.5 rounded-lg text-xs font-semibold cursor-pointer transition-all"
                          style={{ backgroundColor: "#FEE2E2", color: "#991B1B" }}
                        >
                          {deleting === book.id ? "…" : "Sil"}
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
            {filtered.length === 0 && (
              <div className="flex items-center justify-center h-24">
                <p className="text-sm" style={{ color: "var(--muted-foreground)" }}>Kitap bulunamadı</p>
              </div>
            )}
          </div>
        )}
      </div>

      {modal.open && (
        <Modal
          title={modal.editing ? "Kitabı Düzenle" : "Yeni Kitap Ekle"}
          onClose={() => setModal({ open: false, editing: null })}
        >
          <form onSubmit={save} className="space-y-4">
            <FormField label="Başlık *" value={form.title} onChange={(v) => setForm((f) => ({ ...f, title: v }))} required />
            <FormField label="Yazar *" value={form.author} onChange={(v) => setForm((f) => ({ ...f, author: v }))} required />
            <div className="grid grid-cols-2 gap-4">
              <FormField label="ISBN" value={form.isbn ?? ""} onChange={(v) => setForm((f) => ({ ...f, isbn: v }))} />
              <FormField label="Tür" value={form.genre ?? ""} onChange={(v) => setForm((f) => ({ ...f, genre: v }))} />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <FormField label="Yayın Yılı" type="number" value={String(form.published_year ?? "")} onChange={(v) => setForm((f) => ({ ...f, published_year: v ? Number(v) : undefined }))} />
              <div>
                <label className="block text-sm font-semibold mb-1.5" style={{ color: "var(--foreground)" }}>Durum</label>
                <select
                  value={form.available ? "true" : "false"}
                  onChange={(e) => setForm((f) => ({ ...f, available: e.target.value === "true" }))}
                  className="w-full px-4 py-2.5 rounded-xl text-sm outline-none"
                  style={{ backgroundColor: "var(--muted)", border: "1.5px solid var(--border)", color: "var(--foreground)" }}
                >
                  <option value="true">Mevcut</option>
                  <option value="false">Dolu</option>
                </select>
              </div>
            </div>
            <div>
              <label className="block text-sm font-semibold mb-1.5" style={{ color: "var(--foreground)" }}>Açıklama</label>
              <textarea
                value={form.description ?? ""}
                onChange={(e) => setForm((f) => ({ ...f, description: e.target.value }))}
                rows={3}
                className="w-full px-4 py-2.5 rounded-xl text-sm outline-none resize-none"
                style={{ backgroundColor: "var(--muted)", border: "1.5px solid var(--border)", color: "var(--foreground)" }}
              />
            </div>
            <div className="flex gap-3 pt-2">
              <button
                type="button"
                onClick={() => setModal({ open: false, editing: null })}
                className="flex-1 py-2.5 rounded-xl text-sm font-semibold cursor-pointer"
                style={{ backgroundColor: "var(--muted)", color: "var(--muted-foreground)" }}
              >
                İptal
              </button>
              <button
                type="submit"
                disabled={saving}
                className="flex-1 py-2.5 rounded-xl text-sm font-semibold cursor-pointer"
                style={{ backgroundColor: "var(--primary)", color: "var(--primary-foreground)" }}
              >
                {saving ? "Kaydediliyor…" : modal.editing ? "Güncelle" : "Ekle"}
              </button>
            </div>
          </form>
        </Modal>
      )}

      {toast && (
        <div
          className="fixed bottom-6 right-6 px-5 py-3 rounded-xl text-sm font-semibold shadow-lg z-50"
          style={{ backgroundColor: "var(--sidebar)", color: "var(--sidebar-foreground)" }}
        >
          {toast}
        </div>
      )}
    </div>
  );
}

function FormField({ label, value, onChange, required, type = "text" }: { label: string; value: string; onChange: (v: string) => void; required?: boolean; type?: string }) {
  return (
    <div>
      <label className="block text-sm font-semibold mb-1.5" style={{ color: "var(--foreground)" }}>{label}</label>
      <input
        type={type}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        required={required}
        className="w-full px-4 py-2.5 rounded-xl text-sm outline-none"
        style={{ backgroundColor: "var(--muted)", border: "1.5px solid var(--border)", color: "var(--foreground)" }}
      />
    </div>
  );
}

function Modal({ title, onClose, children }: { title: string; onClose: () => void; children: React.ReactNode }) {
  return (
    <div className="fixed inset-0 z-40 flex items-center justify-center p-4" style={{ backgroundColor: "rgba(0,0,0,0.4)" }}>
      <div className="w-full max-w-lg rounded-2xl overflow-hidden shadow-2xl" style={{ backgroundColor: "var(--card)" }}>
        <div className="flex items-center justify-between px-6 py-5" style={{ borderBottom: "1px solid var(--border)" }}>
          <h2 className="font-display text-xl" style={{ color: "var(--foreground)" }}>{title}</h2>
          <button onClick={onClose} className="text-xl leading-none cursor-pointer opacity-40 hover:opacity-100 transition-opacity" style={{ color: "var(--foreground)" }}>✕</button>
        </div>
        <div className="px-6 py-5">{children}</div>
      </div>
    </div>
  );
}

function LoadingSkeleton() {
  return (
    <div className="rounded-2xl overflow-hidden animate-pulse" style={{ backgroundColor: "var(--card)", border: "1px solid var(--border)" }}>
      {Array.from({ length: 6 }).map((_, i) => (
        <div key={i} className="px-5 py-4 flex gap-4" style={{ borderBottom: "1px solid var(--border)" }}>
          <div className="h-4 rounded flex-1" style={{ backgroundColor: "var(--muted)" }} />
          <div className="h-4 rounded w-24" style={{ backgroundColor: "var(--muted)" }} />
          <div className="h-4 rounded w-16" style={{ backgroundColor: "var(--muted)" }} />
        </div>
      ))}
    </div>
  );
}
