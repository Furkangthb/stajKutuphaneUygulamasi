import { useEffect, useState } from "react";
import { api } from "../api";

// Backend'deki Book modelinize tam uyumlu arayüz (api.ts'de de güncellemeyi unutmayın)
export interface Book {
  id: number;
  isbn: string;
  title: string;
  author: string;
  genre: string;
  publish_date: string; // Backend time.Time beklediği için string (ISO formatı) olacak
  description: string;
  available?: boolean;
  reservedUserId?: number; // "Dolu" seçildiğinde kime rezerve edildiği (sadece UI state)
}

const empty: Omit<Book, "id"> = {
  title: "",
  author: "",
  isbn: "",
  genre: "",
  publish_date: "", 
  description: "",
  available: true,
  reservedUserId: undefined,
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
  const [users, setUsers] = useState<{ id: number; first_name: string; last_name: string; email: string }[]>([]);
  const [originalReservedUserId, setOriginalReservedUserId] = useState<number | undefined>(undefined);

  useEffect(() => {
    load();
    api.getUsers(1, 200)
      .then((res: any) => setUsers(Array.isArray(res) ? res : (res?.data || [])))
      .catch(() => {}); // kullanıcı listesi gelmezse "Dolu" seçimi devre dışı kalır, sessiz geç
  }, []);

  const load = () => {
    setLoading(true);
    api.getBooks()
      .then((res: any) => {
        let dataList: any[] = [];
        if (Array.isArray(res)) {
            dataList = res;
        } else if (res && Array.isArray(res.data)) {
            dataList = res.data;
        } else if (res && Array.isArray(res.books)) {
            dataList = res.books;
        }
        setBooks(dataList);
      })
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  };

  const showToast = (msg: string) => {
    setToast(msg);
    setTimeout(() => setToast(""), 3000);
  };

  const openAdd = () => {
    setForm(empty);
    setOriginalReservedUserId(undefined);
    setModal({ open: true, editing: null });
  };

  const openEdit = async (book: Book) => {
    
    if (!book) return;
    
    // Backend'den gelen 2026-08-25T10:00:00Z formatındaki tarihi input type="date" için YYYY-MM-DD'ye çeviriyoruz
    const dateStr = book.publish_date ? book.publish_date.split("T")[0] : "";

    let reservedUserId: number | undefined = undefined;
    if (book.available === false) {
      // Bu kitaba ait aktif rezervasyonu bul, kime ait olduğunu formda göster
      try {
        const all = await api.getAllReservations();
        const active = (Array.isArray(all) ? all : []).find((r) => r.book_id === book.id && r.status === "active");
        reservedUserId = active?.user_id;
      } catch {
        // bulunamazsa boş bırak, admin manuel seçer
      }
    }

    setOriginalReservedUserId(reservedUserId);
    setForm({ 
        title: book.title || "", 
        author: book.author || "", 
        isbn: book.isbn || (book as any).ısbn || "",
        genre: book.genre || "", 
        publish_date: dateStr, 
        available: book.available ?? true, 
        reservedUserId,
        description: book.description || "" 
    });
    setModal({ open: true, editing: book });
  };

  const save = async (e: React.FormEvent) => {
    e.preventDefault();

    if (form.available === false && !form.reservedUserId) {
      showToast("Kitabı 'Dolu' yapmak için kime rezerve edildiğini seçmelisiniz");
      return;
    }

    setSaving(true);
    try {
      // Go (Gin) backend'inin time.Time'ı hatasız parse edebilmesi için tarihi ISO 8601 formatına çeviriyoruz
      const payload = {
          ...form,
          // Eğer tarih seçildiyse sonuna saat ekleyip ISO yap, yoksa şu anki tarihi kullan
          publish_date: form.publish_date ? new Date(form.publish_date).toISOString() : new Date().toISOString()
      };

      let bookId = modal.editing?.id;
      let updatedBookForState: Book;

      if (modal.editing) {
        const updated = await api.updateBook(modal.editing.id, payload as any);
        updatedBookForState = { ...modal.editing, ...updated } as Book;
        setBooks((bs) => (Array.isArray(bs) ? bs : []).map((b) => (b?.id === modal.editing!.id ? updatedBookForState : b)));
      } else {
        const added = await api.addBook(payload as any);
        bookId = added.id;
        updatedBookForState = added;
        setBooks((bs) => [added, ...(Array.isArray(bs) ? bs : [])]);
      }

      // "Dolu" -> önceden farklı biri rezerveliyse ya da yeni Dolu yapıldıysa yeni rezervasyon aç
      const wasAvailable = modal.editing ? modal.editing.available !== false : true;
      const holderChanged = form.reservedUserId !== originalReservedUserId;

      if (form.available === false && bookId && (wasAvailable || holderChanged)) {
        // Önce eski rezervasyon varsa kapat (kişi değiştiyse)
        if (!wasAvailable) {
          await cancelActiveReservation(bookId);
        }
        await api.createReservation({ book_id: bookId, user_id: form.reservedUserId! } as any);
      } else if (form.available !== false && !wasAvailable && bookId) {
        // "Mevcut"a çekildi, önceki aktif rezervasyonu iptal et
        await cancelActiveReservation(bookId);
      }

      showToast(modal.editing ? "Kitap güncellendi" : "Kitap eklendi");
      setModal({ open: false, editing: null });
      load(); // durum (available) her zaman rezervasyon tablosundan hesaplandığı için listeyi tazele
    } catch (e: unknown) {
      showToast(e instanceof Error ? e.message : "İşlem başarısız");
    } finally {
      setSaving(false);
    }
  };

  const cancelActiveReservation = async (bookId: number) => {
    const all = await api.getAllReservations();
    const active = (Array.isArray(all) ? all : []).find((r) => r.book_id === bookId && r.status === "active");
    if (active) await api.updateReservation(active.id, { status: "cancelled" } as any);
  };

  const del = async (book: Book) => {
    if (!book || !confirm(`"${book.title || 'Bu kitabı'}" silmek istediğinizden emin misiniz?`)) return;
    setDeleting(book.id);
    try {
      await api.deleteBook(book.id);
      setBooks((bs) => (Array.isArray(bs) ? bs : []).filter((b) => b?.id !== book.id));
      showToast("Kitap silindi");
    } catch (e: unknown) {
      showToast(e instanceof Error ? e.message : "Silme başarısız");
    } finally {
      setDeleting(null);
    }
  };

  const safeBooks = Array.isArray(books) ? books : [];
  
  const filtered = safeBooks.filter((b) => {
      if (!b) return false;
      const title = b.title || "";
      const author = b.author || "";
      const searchTerm = search || "";
      return title.toLowerCase().includes(searchTerm.toLowerCase()) ||
             author.toLowerCase().includes(searchTerm.toLowerCase());
  });

  return (
    <div className="h-full flex flex-col">
      <div className="px-8 pt-8 pb-6 flex items-start justify-between gap-4">
        <div>
          <h1 className="font-display text-3xl mb-1" style={{ color: "var(--foreground)" }}>
            Kitap Yönetimi
          </h1>
          <p className="text-sm" style={{ color: "var(--muted-foreground)" }}>
            {safeBooks.length} kitap
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
                  {["Başlık", "Yazar", "Tür", "Yayın Tarihi", "Durum", ""].map((h) => (
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
                {filtered.map((book, i) => {
                  if (!book) return null;
                  return (
                  <tr
                    key={book.id || i}
                    style={{ borderBottom: i < filtered.length - 1 ? "1px solid var(--border)" : "none" }}
                    className="hover:bg-[var(--muted)] transition-colors duration-100"
                  >
                    <td className="px-5 py-3 font-semibold" style={{ color: "var(--foreground)" }}>{book.title || "İsimsiz"}</td>
                    <td className="px-5 py-3" style={{ color: "var(--muted-foreground)" }}>{book.author || "—"}</td>
                    <td className="px-5 py-3" style={{ color: "var(--muted-foreground)" }}>{book.genre || "—"}</td>
                    <td className="px-5 py-3" style={{ color: "var(--muted-foreground)" }}>
                      {book.publish_date ? new Date(book.publish_date).toLocaleDateString("tr-TR") : "—"}
                    </td>
                    <td className="px-5 py-3">
                      <span
                        className="px-2 py-0.5 rounded-full text-xs font-bold"
                        style={{ backgroundColor: book.available !== false ? "#D1FAE5" : "#FEE2E2", color: book.available !== false ? "#065F46" : "#991B1B" }}
                      >
                        {book.available !== false ? "Mevcut" : "Dolu"}
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
                )})}
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
              <FormField label="ISBN *" value={form.isbn} onChange={(v) => setForm((f) => ({ ...f, isbn: v }))} required />
              <FormField label="Tür *" value={form.genre} onChange={(v) => setForm((f) => ({ ...f, genre: v }))} required />
            </div>
            
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-semibold mb-1.5" style={{ color: "var(--foreground)" }}>Yayın Tarihi *</label>
                <input
                  type="date"
                  required
                  value={form.publish_date}
                  onChange={(e) => setForm((f) => ({ ...f, publish_date: e.target.value }))}
                  className="w-full px-4 py-2.5 rounded-xl text-sm outline-none"
                  style={{ backgroundColor: "var(--muted)", border: "1.5px solid var(--border)", color: "var(--foreground)" }}
                />
              </div>
              
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

            {/* Yeni Aramalı Seçim (Searchable Dropdown) Bölümü */}
            {form.available === false && (
              <div>
                <label className="block text-sm font-semibold mb-1.5" style={{ color: "var(--foreground)" }}>Kime rezerve edildi? *</label>
                <SearchableUserSelect 
                  users={users} 
                  value={form.reservedUserId} 
                  onChange={(val) => setForm((f) => ({ ...f, reservedUserId: val }))} 
                />
              </div>
            )}
            
            <div>
              <label className="block text-sm font-semibold mb-1.5" style={{ color: "var(--foreground)" }}>Açıklama *</label>
              <textarea
                value={form.description}
                onChange={(e) => setForm((f) => ({ ...f, description: e.target.value }))}
                required
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

// --- YARDIMCI BİLEŞENLER ---

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

function SearchableUserSelect({ users, value, onChange }: { users: any[], value: number | undefined, onChange: (val: number | undefined) => void }) {
  const [search, setSearch] = useState("");
  const [isOpen, setIsOpen] = useState(false);

  const selectedUser = users.find(u => u.id === value);

  // 1. Dışarıdan veya içeriden bir seçim yapıldığında input içine E-POSTA yaz
  useEffect(() => {
    if (selectedUser) {
      setSearch(selectedUser.email);
    } else {
      setSearch("");
    }
  }, [selectedUser]);

  // 2. Arama filtresini E-POSTA'ya göre yap
  const filtered = users.filter(u => 
    (u.email || "").toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="relative">
      <input
        type="text"
        placeholder="E-posta ara ve seç..."
        value={search}
        onFocus={() => setIsOpen(true)}
        onBlur={() => setTimeout(() => setIsOpen(false), 200)}
        onChange={(e) => {
           setSearch(e.target.value);
           setIsOpen(true);
           onChange(undefined);
        }}
        className="w-full px-4 py-2.5 rounded-xl text-sm outline-none"
        style={{ backgroundColor: "var(--muted)", border: "1.5px solid var(--border)", color: "var(--foreground)" }}
      />
      
      {/* Açılır Liste (Dropdown) */}
      {isOpen && (
        <div 
          className="absolute top-full mt-1 left-0 right-0 z-50 rounded-xl shadow-lg max-h-48 overflow-y-auto" 
          style={{ backgroundColor: "var(--card)", border: "1px solid var(--border)" }}
        >
          {filtered.length === 0 ? (
            <div className="px-4 py-3 text-sm opacity-50 text-center">E-posta bulunamadı</div>
          ) : (
            filtered.map(u => (
              <div
                key={u.id}
                className="px-4 py-2.5 text-sm cursor-pointer hover:bg-black/5 dark:hover:bg-white/5 transition-colors flex justify-between items-center"
                style={{ borderBottom: "1px solid var(--border)" }}
                onMouseDown={() => { 
                  onChange(u.id);
                  setSearch(u.email); // Seçildiğinde kutuya e-postayı yaz
                  setIsOpen(false);
                }}
              >
                {/* 3. Listede e-postayı göster (yanında bilgi amaçlı silik renkte isim de yazsın) */}
                <span className="font-medium">{u.email}</span>
                <span className="text-xs opacity-50">{u.first_name} {u.last_name}</span>
              </div>
            ))
          )}
        </div>
      )}
    </div>
  );
}