import { useEffect, useState } from "react";
import { api, type Reservation, type ReservationFull } from "../api";

const STATUS_LABELS: Record<string, string> = {
  active: "Aktif",
  completed: "İade Edildi",
  cancelled: "İptal Edildi",
  expired: "Süresi Doldu",
};

const STATUS_COLORS: Record<string, { bg: string; text: string }> = {
  active: { bg: "#D1FAE5", text: "#065F46" },
  completed: { bg: "#DBEAFE", text: "#1E40AF" },
  cancelled: { bg: "#FEE2E2", text: "#991B1B" },
  expired: { bg: "#FEF3C7", text: "#92400E" },
};

interface ReservationsPageProps {
  userId: number;
}

export default function ReservationsPage({ userId }: ReservationsPageProps) {
  // Veritabanından kitap adı da geleceği için ReservationFull tipine esnetildi
  const [reservations, setReservations] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [updating, setUpdating] = useState<number | null>(null);
  const [toast, setToast] = useState("");

  useEffect(() => {
    api.getUserReservations(userId)
      .then((res: any) => {
        // Güvenli dizi (array) kontrolü
        const dataList = Array.isArray(res) ? res : (res?.data || res?.reservations || []);
        setReservations(dataList);
      })
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, [userId]);

  const showToast = (msg: string) => {
    setToast(msg);
    setTimeout(() => setToast(""), 3000);
  };

  const updateStatus = async (id: number, status: string) => {
    if (!confirm(status === "completed" ? "Kitabı iade ettiğinizi onaylıyor musunuz?" : "Rezervasyonu iptal etmek istiyor musunuz?")) return;
    
    setUpdating(id);
    try {
      await api.updateReservation(id, { status } as any);
      setReservations((rs) => rs.map((r) => (r.id === id ? { ...r, status } : r)));
      showToast("İşlem başarılı");
    } catch (e: unknown) {
      showToast(e instanceof Error ? e.message : "İşlem başarısız");
    } finally {
      setUpdating(null);
    }
  };

  return (
    <div className="h-full flex flex-col">
      <div className="px-8 pt-8 pb-6">
        <h1 className="font-display text-3xl mb-1" style={{ color: "var(--foreground)" }}>
          Rezervasyonlarım
        </h1>
        <p className="text-sm" style={{ color: "var(--muted-foreground)" }}>
          {reservations.length} rezervasyon
        </p>
      </div>

      <div className="flex-1 overflow-y-auto px-8 pb-8">
        {loading ? (
          <LoadingList />
        ) : error ? (
          <div className="flex items-center justify-center h-48">
            <p className="text-sm font-semibold" style={{ color: "#DC2626" }}>{error}</p>
          </div>
        ) : reservations.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-48 gap-2">
            <p className="text-sm font-semibold" style={{ color: "var(--foreground)" }}>Henüz rezervasyonunuz yok</p>
            <p className="text-xs" style={{ color: "var(--muted-foreground)" }}>Kitap kataloğundan rezervasyon yapabilirsiniz</p>
          </div>
        ) : (
          <div className="space-y-3 max-w-2xl">
            {reservations.map((res) => (
              <ReservationCard
                key={res.id}
                reservation={res}
                onUpdate={updateStatus}
                updating={updating === res.id}
              />
            ))}
          </div>
        )}
      </div>

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

function ReservationCard({
  reservation, onUpdate, updating,
}: {
  reservation: any;
  onUpdate: (id: number, status: string) => void;
  updating: boolean;
}) {
  const sc = STATUS_COLORS[reservation.status] ?? { bg: "#F3F4F6", text: "#6B7280" };
  
  // Backend bazen kitap adını göndermeyebilir, defansif kodlama:
  const title = reservation.book_title || reservation.book?.title || `Kitap #${reservation.book_id}`;
  const author = reservation.book_author || reservation.book?.author || "";

  // Gecikme Kontrolü
  const isOverdue = reservation.due_date && new Date(reservation.due_date) < new Date() && reservation.status === "active";

  return (
    <div
      className="rounded-2xl p-5 flex items-start gap-4 transition-all"
      style={{ 
        backgroundColor: "var(--card)", 
        border: `1px solid ${isOverdue ? "#FCA5A5" : "var(--border)"}` 
      }}
    >
      <div
        className="w-12 h-16 rounded-lg flex items-center justify-center shrink-0"
        style={{ backgroundColor: isOverdue ? "#FEF2F2" : "var(--muted)" }}
      >
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke={isOverdue ? "#DC2626" : "var(--muted-foreground)"} strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
          <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z" />
          <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z" />
        </svg>
      </div>
      <div className="flex-1 min-w-0">
        <div className="flex items-start justify-between gap-2 mb-1">
          <h3 className="font-semibold text-sm leading-tight" style={{ color: "var(--foreground)" }}>{title}</h3>
          <span
            className="shrink-0 px-2 py-0.5 rounded-full text-[11px] font-bold uppercase tracking-wider"
            style={{ backgroundColor: sc.bg, color: sc.text }}
          >
            {STATUS_LABELS[reservation.status] || "Bilinmiyor"}
          </span>
        </div>
        {author && <p className="text-xs mb-2" style={{ color: "var(--muted-foreground)" }}>{author}</p>}
        
        <div className="flex gap-4 mt-1">
          {reservation.reserved_at && (
            <p className="text-xs" style={{ color: "var(--muted-foreground)" }}>
              <span className="opacity-70">Alış:</span> {new Date(reservation.reserved_at).toLocaleDateString("tr-TR")}
            </p>
          )}
          {reservation.due_date && (
            <p className="text-xs font-medium" style={{ color: isOverdue ? "#DC2626" : "var(--muted-foreground)" }}>
              <span className="opacity-70">İade Tarihi:</span> {new Date(reservation.due_date).toLocaleDateString("tr-TR")}
              {isOverdue && <span className="ml-1.5 bg-red-100 text-red-700 px-1 py-0.5 rounded text-[9px] uppercase tracking-wider font-bold">Gecikti</span>}
            </p>
          )}
        </div>

        {/* Butonlar: Sadece aktif rezervasyonlarda görünür */}
        {reservation.status === "active" && (
          <div className="flex gap-2 mt-4">
            <button
              disabled={updating}
              onClick={() => onUpdate(reservation.id, "completed")}
              className="px-4 py-2 rounded-lg text-xs font-semibold cursor-pointer transition-all"
              style={{ backgroundColor: "var(--primary)", color: "var(--primary-foreground)" }}
            >
              {updating ? "İşleniyor…" : "Teslim Ettim"}
            </button>
            <button
              disabled={updating}
              onClick={() => onUpdate(reservation.id, "cancelled")}
              className="px-4 py-2 rounded-lg text-xs font-semibold cursor-pointer transition-all"
              style={{ backgroundColor: "#FEE2E2", color: "#991B1B" }}
            >
              {updating ? "İşleniyor…" : "İptal Et"}
            </button>
          </div>
        )}
      </div>
    </div>
  );
}

function LoadingList() {
  return (
    <div className="space-y-3 max-w-2xl">
      {Array.from({ length: 4 }).map((_, i) => (
        <div
          key={i}
          className="rounded-2xl p-5 flex gap-4 animate-pulse"
          style={{ backgroundColor: "var(--card)", border: "1px solid var(--border)" }}
        >
          <div className="w-12 h-16 rounded-lg" style={{ backgroundColor: "var(--muted)" }} />
          <div className="flex-1 space-y-2">
            <div className="h-4 rounded" style={{ backgroundColor: "var(--muted)", width: "60%" }} />
            <div className="h-3 rounded" style={{ backgroundColor: "var(--muted)", width: "40%" }} />
          </div>
        </div>
      ))}
    </div>
  );
}