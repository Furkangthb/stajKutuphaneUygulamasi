import { useEffect, useState } from "react";
import { api, type Reservation, type ReservationFull } from "../api";

const STATUS_LABELS: Record<string, string> = {
  pending: "Onay Bekliyor",
  active: "Aktif",
  completed: "İade Edildi",
  cancelled: "İptal Edildi",
  expired: "Süresi Doldu",
};

const STATUS_COLORS: Record<string, { bg: string; text: string }> = {
  pending: { bg: "#EDE9FE", text: "#5B21B6" },
  active: { bg: "#D1FAE5", text: "#065F46" },
  completed: { bg: "#DBEAFE", text: "#1E40AF" },
  cancelled: { bg: "#FEE2E2", text: "#991B1B" },
  expired: { bg: "#FEF3C7", text: "#92400E" },
};

const STATUS_FILTERS = ["Tümü", "Onay Bekliyor", "Aktif", "İade Edildi", "İptal Edildi", "Süresi Doldu"];
const STATUS_MAP: Record<string, string> = {
  "Onay Bekliyor": "pending", "Aktif": "active", "İade Edildi": "completed", "İptal Edildi": "cancelled", "Süresi Doldu": "expired",
};

interface ReservationsPageProps {
  userId: number;
}

export default function ReservationsPage({ userId }: ReservationsPageProps) {
  const [reservations, setReservations] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [statusFilter, setStatusFilter] = useState("Tümü");
  const [cancelingId, setCancelingId] = useState<number | null>(null);
  const [toast, setToast] = useState("");

  useEffect(() => {
    api.getUserReservations(userId)
      .then((res: any) => {
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

  const cancelReservation = async (id: number) => {
    setCancelingId(id);
    try {
      await api.updateReservation(id, { status: "cancelled" });
      setReservations((rs) => rs.map((r) => (r.id === id ? { ...r, status: "cancelled" } : r)));
      showToast("Rezervasyon iptal edildi");
    } catch (e: unknown) {
      showToast(e instanceof Error ? e.message : "İptal edilemedi");
    } finally {
      setCancelingId(null);
    }
  };

  const filteredReservations = reservations.filter(
    (r) => statusFilter === "Tümü" || r.status === STATUS_MAP[statusFilter]
  );

  return (
    <div className="h-full flex flex-col">
      <div className="px-8 pt-8 pb-6">
        <h1 className="font-display text-3xl mb-1" style={{ color: "var(--foreground)" }}>
          Rezervasyonlarım
        </h1>
        <p className="text-sm" style={{ color: "var(--muted-foreground)" }}>
          {filteredReservations.length} rezervasyon
        </p>
      </div>

      <div className="px-8 pb-5 flex gap-1.5 flex-wrap">
        {STATUS_FILTERS.map((s) => (
          <button
            key={s}
            onClick={() => setStatusFilter(s)}
            className="px-3 py-2 rounded-lg text-xs font-semibold transition-all cursor-pointer"
            style={{
              backgroundColor: statusFilter === s ? "var(--primary)" : "var(--card)",
              color: statusFilter === s ? "var(--primary-foreground)" : "var(--muted-foreground)",
              border: `1.5px solid ${statusFilter === s ? "var(--primary)" : "var(--border)"}`,
            }}
          >
            {s}
          </button>
        ))}
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
        ) : filteredReservations.length === 0 ? (
          <div className="flex items-center justify-center h-48">
            <p className="text-sm" style={{ color: "var(--muted-foreground)" }}>Bu durumda rezervasyon bulunamadı</p>
          </div>
        ) : (
          <div className="space-y-3 max-w-2xl">
            {filteredReservations.map((res) => (
              <ReservationCard
                key={res.id}
                reservation={res}
                canceling={cancelingId === res.id}
                onCancel={() => cancelReservation(res.id)}
              />
            ))}
          </div>
        )}
      </div>

      {toast && (
        <div className="fixed bottom-6 right-6 px-5 py-3 rounded-xl text-sm font-semibold shadow-lg z-50" style={{ backgroundColor: "var(--sidebar)", color: "var(--sidebar-foreground)" }}>
          {toast}
        </div>
      )}
    </div>
  );
}

function ReservationCard({ reservation, canceling, onCancel }: { reservation: any; canceling: boolean; onCancel: () => void }) {
  const sc = STATUS_COLORS[reservation.status] ?? { bg: "#F3F4F6", text: "#6B7280" };

  const title = reservation.book_title || reservation.book?.title || `Kitap #${reservation.book_id}`;
  const author = reservation.book_author || reservation.book?.author || "";

  const isOverdue = reservation.due_date && new Date(reservation.due_date) < new Date() && reservation.status === "active";
  // Sadece onay bekleyen (henüz kitap teslim edilmemiş) talepler kullanıcı tarafından iptal edilebilir.
  // "active" (kitap elinde) durumundaki bir rezervasyon fiziksel iade gerektirdiği için sadece kütüphaneci "tamamlandı" yapabilir.
  const canCancel = reservation.status === "pending";

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

        {canCancel && (
          <div className="mt-3">
            <button
              onClick={onCancel}
              disabled={canceling}
              className="px-3 py-1.5 rounded-lg text-xs font-bold cursor-pointer"
              style={{ backgroundColor: "#FEE2E2", color: "#991B1B" }}
            >
              {canceling ? "İptal ediliyor…" : "Talebi İptal Et"}
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