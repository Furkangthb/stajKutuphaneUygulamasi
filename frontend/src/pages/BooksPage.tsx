import { useEffect, useState } from "react";
import { api, type Book,type  Reservation } from "../api";

interface BooksPageProps {
  userId: number;
}

const GENRES = ["Tümü", "Roman", "Bilim", "Tarih", "Felsefe", "Teknoloji", "Biyografi", "Şiir"];

export default function BooksPage({ userId }: BooksPageProps) {
  const [books, setBooks] = useState<Book[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [search, setSearch] = useState("");
  const [genre, setGenre] = useState("Tümü");
  const [reserving, setReserving] = useState<number | null>(null);
  const [toast, setToast] = useState("");

  useEffect(() => {
    api.getBooks()
      .then(setBooks)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, []);

  const filtered = books.filter((b) => {
    const matchSearch =
      b.title.toLowerCase().includes(search.toLowerCase()) ||
      b.author.toLowerCase().includes(search.toLowerCase());
    const matchGenre = genre === "Tümü" || b.genre === genre;
    return matchSearch && matchGenre;
  });

  const showToast = (msg: string) => {
    setToast(msg);
    setTimeout(() => setToast(""), 3000);
  };

  const reserve = async (book: Book) => {
    setReserving(book.id);
    try {
      await api.createReservation({ book_id: book.id });
      showToast(`"${book.title}" başarıyla rezerve edildi!`);
      setBooks((bs) => bs.map((b) => (b.id === book.id ? { ...b, available: false } : b)));
    } catch (e: unknown) {
      showToast(e instanceof Error ? e.message : "Rezervasyon başarısız");
    } finally {
      setReserving(null);
    }
  };

  return (
    <div className="h-full flex flex-col">
      {/* Header */}
      <div className="px-8 pt-8 pb-6">
        <h1 className="font-display text-3xl mb-1" style={{ color: "var(--foreground)" }}>
          Kitap Kataloğu
        </h1>
        <p className="text-sm" style={{ color: "var(--muted-foreground)" }}>
          {books.length} kitap mevcut
        </p>
      </div>

      {/* Filters */}
      <div className="px-8 pb-4 flex flex-wrap gap-3">
        <div className="relative flex-1 min-w-48">
          <SearchIcon />
          <input
            type="text"
            placeholder="Kitap veya yazar ara…"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full pl-9 pr-4 py-2.5 rounded-xl text-sm outline-none"
            style={{
              backgroundColor: "var(--card)",
              border: "1.5px solid var(--border)",
              color: "var(--foreground)",
            }}
          />
        </div>
        <div className="flex gap-1.5 flex-wrap">
          {GENRES.map((g) => (
            <button
              key={g}
              onClick={() => setGenre(g)}
              className="px-3 py-2 rounded-lg text-xs font-semibold transition-all cursor-pointer"
              style={{
                backgroundColor: genre === g ? "var(--primary)" : "var(--card)",
                color: genre === g ? "var(--primary-foreground)" : "var(--muted-foreground)",
                border: `1.5px solid ${genre === g ? "var(--primary)" : "var(--border)"}`,
              }}
            >
              {g}
            </button>
          ))}
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto px-8 pb-8">
        {loading ? (
          <LoadingGrid />
        ) : error ? (
          <ErrorState message={error} />
        ) : filtered.length === 0 ? (
          <EmptyState />
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-4">
            {filtered.map((book) => (
              <BookCard
                key={book.id}
                book={book}
                onReserve={reserve}
                reserving={reserving === book.id}
              />
            ))}
          </div>
        )}
      </div>

      {toast && (
        <div
          className="fixed bottom-6 right-6 px-5 py-3 rounded-xl text-sm font-semibold shadow-lg z-50 transition-all"
          style={{ backgroundColor: "var(--sidebar)", color: "var(--sidebar-foreground)" }}
        >
          {toast}
        </div>
      )}
    </div>
  );
}

function BookCard({ book, onReserve, reserving }: { book: Book; onReserve: (b: Book) => void; reserving: boolean }) {
  const colors = ["#DBEAFE", "#FEF3C7", "#D1FAE5", "#F3E8FF", "#FFE4E6"];
  const bg = colors[book.id % colors.length];

  return (
    <div
      className="rounded-2xl overflow-hidden transition-all duration-200 hover:-translate-y-0.5"
      style={{ backgroundColor: "var(--card)", border: "1px solid var(--border)", boxShadow: "0 1px 3px rgba(0,0,0,0.05)" }}
    >
      <div className="h-36 flex items-center justify-center" style={{ backgroundColor: bg }}>
        <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
          <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z" />
          <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z" />
        </svg>
      </div>
      <div className="p-5">
        <div className="flex items-start justify-between gap-2 mb-1">
          <h3 className="font-display text-base leading-snug line-clamp-2" style={{ color: "var(--foreground)" }}>
            {book.title}
          </h3>
          <span
            className="shrink-0 mt-0.5 px-2 py-0.5 rounded-full text-xs font-bold"
            style={{
              backgroundColor: book.available ? "#D1FAE5" : "#FEE2E2",
              color: book.available ? "#065F46" : "#991B1B",
            }}
          >
            {book.available ? "Mevcut" : "Dolu"}
          </span>
        </div>
        <p className="text-sm mb-1" style={{ color: "var(--muted-foreground)" }}>{book.author}</p>
        {book.genre && (
          <p className="text-xs mb-4" style={{ color: "var(--muted-foreground)" }}>{book.genre} {book.published_year ? `· ${book.published_year}` : ""}</p>
        )}
        <button
          disabled={!book.available || reserving}
          onClick={() => onReserve(book)}
          className="w-full py-2 rounded-lg text-sm font-semibold transition-all duration-150 cursor-pointer"
          style={{
            backgroundColor: book.available && !reserving ? "var(--primary)" : "var(--muted)",
            color: book.available && !reserving ? "var(--primary-foreground)" : "var(--muted-foreground)",
          }}
        >
          {reserving ? "Rezerve ediliyor…" : book.available ? "Rezerve Et" : "Müsait Değil"}
        </button>
      </div>
    </div>
  );
}

function SearchIcon() {
  return (
    <svg className="absolute left-3 top-1/2 -translate-y-1/2 pointer-events-none" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="var(--muted-foreground)" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <circle cx="11" cy="11" r="8" />
      <line x1="21" y1="21" x2="16.65" y2="16.65" />
    </svg>
  );
}

function LoadingGrid() {
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-4">
      {Array.from({ length: 6 }).map((_, i) => (
        <div key={i} className="rounded-2xl overflow-hidden animate-pulse" style={{ backgroundColor: "var(--card)", border: "1px solid var(--border)" }}>
          <div className="h-36" style={{ backgroundColor: "var(--muted)" }} />
          <div className="p-5 space-y-2">
            <div className="h-4 rounded" style={{ backgroundColor: "var(--muted)", width: "70%" }} />
            <div className="h-3 rounded" style={{ backgroundColor: "var(--muted)", width: "50%" }} />
            <div className="h-8 rounded-lg mt-4" style={{ backgroundColor: "var(--muted)" }} />
          </div>
        </div>
      ))}
    </div>
  );
}

function ErrorState({ message }: { message: string }) {
  return (
    <div className="flex flex-col items-center justify-center h-48 gap-2">
      <p className="text-sm font-semibold" style={{ color: "#DC2626" }}>Hata oluştu</p>
      <p className="text-xs" style={{ color: "var(--muted-foreground)" }}>{message}</p>
    </div>
  );
}

function EmptyState() {
  return (
    <div className="flex flex-col items-center justify-center h-48 gap-2">
      <p className="text-sm font-semibold" style={{ color: "var(--foreground)" }}>Kitap bulunamadı</p>
      <p className="text-xs" style={{ color: "var(--muted-foreground)" }}>Arama kriterlerinizi değiştirin</p>
    </div>
  );
}
