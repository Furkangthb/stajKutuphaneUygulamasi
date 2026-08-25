import { useEffect, useState } from "react";
import { api } from "../api";

// Backend'den gelen User modelinize uyumlu arayüz (api.ts'de de böyle olmalı)
export interface User {
    id: number;
    first_name: string;
    last_name: string;
    email: string;
    role: string;
    phone?: string;
    created_at?: string;
}

export default function ManageUsersPage() {
    const [users, setUsers] = useState<User[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");
    const [deleting, setDeleting] = useState<number | null>(null);
    const [editing, setEditing] = useState<User | null>(null);
    // Backend'deki UserRework struct'ına (first_name, last_name, email, role) uyarlandı
    const [editForm, setEditForm] = useState({ first_name: "", last_name: "", email: "", role: "user" as "user" | "admin" });
    const [saving, setSaving] = useState(false);
    const [toast, setToast] = useState("");
    const [search, setSearch] = useState("");

    useEffect(() => {
        // api.getUsers() artık {data, page} zarfını doğru tipiyle döndürüyor
        api.getUsers()
            .then((res) => {
                setUsers(res.data || []);
            })
            .catch((e) => setError(e.message))
            .finally(() => setLoading(false));
    }, []);

    const showToast = (msg: string) => {
        setToast(msg);
        setTimeout(() => setToast(""), 3000);
    };

    const openEdit = (user: User) => {
        setEditing(user);
        setEditForm({
            first_name: user.first_name || "",
            last_name: user.last_name || "",
            email: user.email || "",
            role: user.role === "admin" ? "admin" : "user"
        });
    };

    const saveEdit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!editing) return;
        setSaving(true);
        try {
            const updated = await api.updateUser(editing.id, editForm);

            // ÇÖZÜM BURADA:
            setUsers((us) => (us || []).map((u) => {
                if (u.id === editing.id) {
                    return {
                        ...u,
                        ...updated,
                        created_at: (!updated.created_at || updated.created_at.startsWith("0001"))
                            ? u.created_at
                            : updated.created_at
                    };
                }
                return u;
            }));

            showToast("Kullanıcı güncellendi");
            setEditing(null);
        } catch (err: unknown) {
            showToast(err instanceof Error ? err.message : "Güncelleme başarısız");
        } finally {
            setSaving(false);
        }
    };

    const del = async (user: User) => {
        const fullName = `${user.first_name || ""} ${user.last_name || ""}`.trim();
        if (!confirm(`"${fullName || 'Bu kullanıcı'}" kullanıcısını silmek istediğinizden emin misiniz?`)) return;
        setDeleting(user.id);
        try {
            await api.deleteUser(user.id);
            setUsers((us) => (us || []).filter((u) => u.id !== user.id));
            showToast("Kullanıcı silindi");
        } catch (err: unknown) {
            showToast(err instanceof Error ? err.message : "Silme başarısız");
        } finally {
            setDeleting(null);
        }
    };

    // Arama fonksiyonu artık first_name ve last_name üzerinde arama yapıyor
    const filtered = (users || []).filter(
        (u) =>
            (u.first_name || "").toLowerCase().includes(search.toLowerCase()) ||
            (u.last_name || "").toLowerCase().includes(search.toLowerCase()) ||
            (u.email || "").toLowerCase().includes(search.toLowerCase())
    );

    return (
        <div className="h-full flex flex-col">
            <div className="px-8 pt-8 pb-6">
                <h1 className="font-display text-3xl mb-1" style={{ color: "var(--foreground)" }}>
                    Kullanıcı Yönetimi
                </h1>
                <p className="text-sm" style={{ color: "var(--muted-foreground)" }}>
                    {(users || []).length} kayıtlı kullanıcı
                </p>
            </div>

            <div className="px-8 pb-4">
                <div className="relative max-w-sm">
                    <svg className="absolute left-3 top-1/2 -translate-y-1/2 pointer-events-none" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="var(--muted-foreground)" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                        <circle cx="11" cy="11" r="8" /><line x1="21" y1="21" x2="16.65" y2="16.65" />
                    </svg>
                    <input
                        type="text"
                        placeholder="İsim veya e-posta ara…"
                        value={search}
                        onChange={(e) => setSearch(e.target.value)}
                        className="w-full pl-9 pr-4 py-2.5 rounded-xl text-sm outline-none"
                        style={{ backgroundColor: "var(--card)", border: "1.5px solid var(--border)", color: "var(--foreground)" }}
                    />
                </div>
            </div>

            <div className="flex-1 overflow-y-auto px-8 pb-8">
                {loading ? (
                    <div className="space-y-2">
                        {Array.from({ length: 5 }).map((_, i) => (
                            <div key={i} className="h-14 rounded-xl animate-pulse" style={{ backgroundColor: "var(--card)" }} />
                        ))}
                    </div>
                ) : error ? (
                    <p className="text-sm" style={{ color: "#DC2626" }}>{error}</p>
                ) : (
                    <div className="rounded-2xl overflow-hidden" style={{ backgroundColor: "var(--card)", border: "1px solid var(--border)" }}>
                        <table className="w-full text-sm">
                            <thead>
                                <tr style={{ borderBottom: "1px solid var(--border)" }}>
                                    {["Kullanıcı", "E-posta", "Rol", "Kayıt Tarihi", ""].map((h) => (
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
                                {filtered.map((user, i) => (
                                    <tr
                                        key={user.id || i}
                                        style={{ borderBottom: i < filtered.length - 1 ? "1px solid var(--border)" : "none" }}
                                        className="transition-colors duration-100"
                                    >
                                        <td className="px-5 py-3">
                                            <div className="flex items-center gap-3">
                                                <div
                                                    className="w-8 h-8 rounded-full flex items-center justify-center text-xs font-bold uppercase"
                                                    style={{ backgroundColor: "var(--secondary)", color: "var(--secondary-foreground)" }}
                                                >
                                                    {user.first_name ? user.first_name[0] : "?"}
                                                </div>
                                                <span className="font-semibold" style={{ color: "var(--foreground)" }}>
                                                    {`${user.first_name || ""} ${user.last_name || ""}`.trim() || "Bilinmiyor"}
                                                </span>
                                            </div>
                                        </td>
                                        <td className="px-5 py-3" style={{ color: "var(--muted-foreground)" }}>{user.email || "—"}</td>
                                        <td className="px-5 py-3">
                                            <span
                                                className="px-2 py-0.5 rounded-full text-xs font-bold capitalize"
                                                style={{
                                                    backgroundColor: user.role === "admin" ? "#FEF3C7" : "#DBEAFE",
                                                    color: user.role === "admin" ? "#92400E" : "#1E40AF",
                                                }}
                                            >
                                                {user.role === "admin" ? "Admin" : "Kullanıcı"}
                                            </span>
                                        </td>
                                        <td className="px-5 py-3" style={{ color: "var(--muted-foreground)" }}>
                                            {user.created_at ? new Date(user.created_at).toLocaleDateString("tr-TR") : "—"}
                                        </td>
                                        <td className="px-5 py-3">
                                            <div className="flex items-center gap-2">
                                                <button
                                                    onClick={() => openEdit(user)}
                                                    className="px-3 py-1.5 rounded-lg text-xs font-semibold cursor-pointer transition-all"
                                                    style={{ backgroundColor: "var(--secondary)", color: "var(--secondary-foreground)" }}
                                                >
                                                    Düzenle
                                                </button>
                                                <button
                                                    onClick={() => del(user)}
                                                    disabled={deleting === user.id}
                                                    className="px-3 py-1.5 rounded-lg text-xs font-semibold cursor-pointer transition-all"
                                                    style={{ backgroundColor: "#FEE2E2", color: "#991B1B" }}
                                                >
                                                    {deleting === user.id ? "…" : "Sil"}
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                        {filtered.length === 0 && (
                            <div className="flex items-center justify-center h-24">
                                <p className="text-sm" style={{ color: "var(--muted-foreground)" }}>Kullanıcı bulunamadı</p>
                            </div>
                        )}
                    </div>
                )}
            </div>

            {editing && (
                <div className="fixed inset-0 z-40 flex items-center justify-center p-4" style={{ backgroundColor: "rgba(0,0,0,0.4)" }}>
                    <div className="w-full max-w-md rounded-2xl overflow-hidden shadow-2xl" style={{ backgroundColor: "var(--card)" }}>
                        <div className="flex items-center justify-between px-6 py-5" style={{ borderBottom: "1px solid var(--border)" }}>
                            <h2 className="font-display text-xl" style={{ color: "var(--foreground)" }}>Kullanıcıyı Düzenle</h2>
                            <button onClick={() => setEditing(null)} className="text-xl leading-none cursor-pointer opacity-40 hover:opacity-100 transition-opacity" style={{ color: "var(--foreground)" }}>✕</button>
                        </div>
                        <form onSubmit={saveEdit} className="px-6 py-5 space-y-4">

                            {/* Form inputları backend struct'ına göre uyarlandı */}
                            <div className="flex gap-4">
                                <div className="flex-1">
                                    <label className="block text-sm font-semibold mb-1.5" style={{ color: "var(--foreground)" }}>Ad</label>
                                    <input
                                        type="text"
                                        value={editForm.first_name}
                                        onChange={(e) => setEditForm((f) => ({ ...f, first_name: e.target.value }))}
                                        required
                                        className="w-full px-4 py-2.5 rounded-xl text-sm outline-none"
                                        style={{ backgroundColor: "var(--muted)", border: "1.5px solid var(--border)", color: "var(--foreground)" }}
                                    />
                                </div>
                                <div className="flex-1">
                                    <label className="block text-sm font-semibold mb-1.5" style={{ color: "var(--foreground)" }}>Soyad</label>
                                    <input
                                        type="text"
                                        value={editForm.last_name}
                                        onChange={(e) => setEditForm((f) => ({ ...f, last_name: e.target.value }))}
                                        required
                                        className="w-full px-4 py-2.5 rounded-xl text-sm outline-none"
                                        style={{ backgroundColor: "var(--muted)", border: "1.5px solid var(--border)", color: "var(--foreground)" }}
                                    />
                                </div>
                            </div>

                            <div>
                                <label className="block text-sm font-semibold mb-1.5" style={{ color: "var(--foreground)" }}>E-posta</label>
                                <input
                                    type="email"
                                    value={editForm.email}
                                    onChange={(e) => setEditForm((f) => ({ ...f, email: e.target.value }))}
                                    required
                                    className="w-full px-4 py-2.5 rounded-xl text-sm outline-none"
                                    style={{ backgroundColor: "var(--muted)", border: "1.5px solid var(--border)", color: "var(--foreground)" }}
                                />
                            </div>

                            <div>
                                <label className="block text-sm font-semibold mb-1.5" style={{ color: "var(--foreground)" }}>Rol</label>
                                <select
                                    value={editForm.role}
                                    onChange={(e) => setEditForm((f) => ({ ...f, role: e.target.value as "user" | "admin" }))}
                                    className="w-full px-4 py-2.5 rounded-xl text-sm outline-none"
                                    style={{ backgroundColor: "var(--muted)", border: "1.5px solid var(--border)", color: "var(--foreground)" }}
                                >
                                    <option value="user">Kullanıcı</option>
                                    <option value="admin">Admin</option>
                                </select>
                            </div>
                            <div className="flex gap-3 pt-2">
                                <button type="button" onClick={() => setEditing(null)} className="flex-1 py-2.5 rounded-xl text-sm font-semibold cursor-pointer" style={{ backgroundColor: "var(--muted)", color: "var(--muted-foreground)" }}>İptal</button>
                                <button type="submit" disabled={saving} className="flex-1 py-2.5 rounded-xl text-sm font-semibold cursor-pointer" style={{ backgroundColor: "var(--primary)", color: "var(--primary-foreground)" }}>
                                    {saving ? "Kaydediliyor…" : "Güncelle"}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
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
