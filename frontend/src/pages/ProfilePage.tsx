import { useState } from "react";
import { api, type User } from "../api";

interface ProfilePageProps {
  user: User;
  onUpdate: (u: User) => void;
}

export default function ProfilePage({ user, onUpdate }: ProfilePageProps) {
  const [form, setForm] = useState({
    first_name: user.first_name ?? "",
    last_name: user.last_name ?? "",
    email: user.email ?? "",
  });
  const [pwForm, setPwForm] = useState({ current: "", next: "", confirm: "" });
  const [saving, setSaving] = useState(false);
  const [savingPw, setSavingPw] = useState(false);
  const [toast, setToast] = useState<{ msg: string; ok: boolean } | null>(null);

  const showToast = (msg: string, ok = true) => {
    setToast({ msg, ok });
    setTimeout(() => setToast(null), 3000);
  };

  const saveProfile = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      await api.updateUser(user.id, {
        first_name: form.first_name,
        last_name: form.last_name,
        email: form.email,
        role: user.role,
      });
      onUpdate({ ...user, first_name: form.first_name, last_name: form.last_name, email: form.email });
      showToast("Profil güncellendi");
    } catch (err: unknown) {
      showToast(err instanceof Error ? err.message : "Güncelleme başarısız", false);
    } finally {
      setSaving(false);
    }
  };

  const savePassword = async (e: React.FormEvent) => {
    e.preventDefault();
    if (pwForm.next !== pwForm.confirm) {
      showToast("Şifreler eşleşmiyor", false);
      return;
    }
    if (pwForm.next.length < 6) {
      showToast("Şifre en az 6 karakter olmalı", false);
      return;
    }
    setSavingPw(true);
    try {
      await api.changePassword(user.id, {
        current_password: pwForm.current,
        new_password: pwForm.next,
      });
      setPwForm({ current: "", next: "", confirm: "" });
      showToast("Şifre güncellendi");
    } catch (err: unknown) {
      showToast(err instanceof Error ? err.message : "Şifre güncellenemedi", false);
    } finally {
      setSavingPw(false);
    }
  };

  const displayName = [user.first_name, user.last_name].filter(Boolean).join(" ") || user.email || "Kullanıcı";
  const initials = (
    (user.first_name?.[0] ?? "") + (user.last_name?.[0] ?? user.first_name?.[1] ?? "")
  ).toUpperCase() || "U";
  const avatarColors = ["#DBEAFE", "#D1FAE5", "#FEF3C7", "#F3E8FF", "#FFE4E6"];
  const avatarBg = avatarColors[user.id % avatarColors.length];

  return (
    <div className="h-full overflow-y-auto">
      <div className="max-w-2xl mx-auto px-8 py-10">

        {/* Avatar card */}
        <div
          className="rounded-2xl p-8 mb-6 flex items-center gap-6"
          style={{ backgroundColor: "var(--card)", border: "1px solid var(--border)" }}
        >
          <div
            className="w-20 h-20 rounded-2xl flex items-center justify-center font-display text-3xl shrink-0"
            style={{ backgroundColor: avatarBg, color: "#1A1917" }}
          >
            {initials}
          </div>
          <div className="flex-1 min-w-0">
            <h2 className="font-display text-2xl mb-0.5" style={{ color: "var(--foreground)" }}>
              {displayName}
            </h2>
            <p className="text-sm mb-3" style={{ color: "var(--muted-foreground)" }}>{user.email}</p>
            <span
              className="px-3 py-1 rounded-full text-xs font-bold capitalize"
              style={{
                backgroundColor: user.role === "admin" ? "#FEF3C7" : "#DBEAFE",
                color: user.role === "admin" ? "#92400E" : "#1E40AF",
              }}
            >
              {user.role === "admin" ? "Admin" : "Kullanıcı"}
            </span>
          </div>
         
        </div>

        {/* Profile form */}
        <Section title="Hesap Bilgileri" subtitle="Kullanıcı adı ve e-posta adresinizi güncelleyin">
          <form onSubmit={saveProfile} className="space-y-4">
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <Field
                label="Ad"
                type="text"
                value={form.first_name}
                onChange={(v) => setForm((f) => ({ ...f, first_name: v }))}
                required
              />
              <Field
                label="Soyad"
                type="text"
                value={form.last_name}
                onChange={(v) => setForm((f) => ({ ...f, last_name: v }))}
                required
              />
            </div>
            <Field
              label="E-posta"
              type="email"
              value={form.email}
              onChange={(v) => setForm((f) => ({ ...f, email: v }))}
              required
            />
            <div className="flex justify-end pt-1">
              <button
                type="submit"
                disabled={saving}
                className="px-5 py-2.5 rounded-xl text-sm font-semibold cursor-pointer transition-all"
                style={{ backgroundColor: "var(--primary)", color: "var(--primary-foreground)" }}
              >
                {saving ? "Kaydediliyor…" : "Kaydet"}
              </button>
            </div>
          </form>
        </Section>

        {/* Password form */}
        <Section title="Şifre Değiştir" subtitle="Hesabınızı güvende tutmak için güçlü bir şifre kullanın">
          <form onSubmit={savePassword} className="space-y-4">
            <Field
              label="Mevcut Şifre"
              type="password"
              value={pwForm.current}
              onChange={(v) => setPwForm((f) => ({ ...f, current: v }))}
              required
            />
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <Field
                label="Yeni Şifre"
                type="password"
                value={pwForm.next}
                onChange={(v) => setPwForm((f) => ({ ...f, next: v }))}
                required
              />
              <Field
                label="Şifre Tekrar"
                type="password"
                value={pwForm.confirm}
                onChange={(v) => setPwForm((f) => ({ ...f, confirm: v }))}
                required
              />
            </div>
            {pwForm.next && pwForm.confirm && pwForm.next !== pwForm.confirm && (
              <p className="text-xs font-medium" style={{ color: "#DC2626" }}>Şifreler eşleşmiyor</p>
            )}
            <div className="flex justify-end pt-1">
              <button
                type="submit"
                disabled={savingPw}
                className="px-5 py-2.5 rounded-xl text-sm font-semibold cursor-pointer transition-all"
                style={{ backgroundColor: "var(--primary)", color: "var(--primary-foreground)" }}
              >
                {savingPw ? "Güncelleniyor…" : "Şifreyi Güncelle"}
              </button>
            </div>
          </form>
        </Section>

        {/* Info card */}
        <div
          className="rounded-2xl p-5 flex items-start gap-4"
          style={{ backgroundColor: "var(--secondary)", border: "1px solid var(--border)" }}
        >
          <div className="mt-0.5">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="var(--secondary-foreground)" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <circle cx="12" cy="12" r="10" />
              <line x1="12" y1="8" x2="12" y2="12" />
              <line x1="12" y1="16" x2="12.01" y2="16" />
            </svg>
          </div>
          <p className="text-xs leading-relaxed" style={{ color: "var(--secondary-foreground)" }}>
            Hesap bilgilerinizde yapılan değişiklikler anında geçerli olur. Şifrenizi unutursanız sistem yöneticinizle iletişime geçin.
          </p>
        </div>

      </div>

      {toast && (
        <div
          className="fixed bottom-6 right-6 px-5 py-3 rounded-xl text-sm font-semibold shadow-lg z-50 transition-all"
          style={{
            backgroundColor: toast.ok ? "var(--sidebar)" : "#DC2626",
            color: "white",
          }}
        >
          {toast.msg}
        </div>
      )}
    </div>
  );
}

function Section({ title, subtitle, children }: { title: string; subtitle: string; children: React.ReactNode }) {
  return (
    <div
      className="rounded-2xl p-6 mb-4"
      style={{ backgroundColor: "var(--card)", border: "1px solid var(--border)" }}
    >
      <div className="mb-5" style={{ borderBottom: "1px solid var(--border)", paddingBottom: "1rem" }}>
        <h3 className="font-display text-lg" style={{ color: "var(--foreground)" }}>{title}</h3>
        <p className="text-xs mt-0.5" style={{ color: "var(--muted-foreground)" }}>{subtitle}</p>
      </div>
      {children}
    </div>
  );
}

function Field({ label, type, value, onChange, required }: {
  label: string; type: string; value: string;
  onChange: (v: string) => void; required?: boolean;
}) {
  return (
    <div>
      <label className="block text-sm font-semibold mb-1.5" style={{ color: "var(--foreground)" }}>{label}</label>
      <input
        type={type}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        required={required}
        className="w-full px-4 py-2.5 rounded-xl text-sm outline-none transition-all"
        style={{ backgroundColor: "var(--muted)", border: "1.5px solid var(--border)", color: "var(--foreground)" }}
        onFocus={(e) => (e.currentTarget.style.borderColor = "var(--primary)")}
        onBlur={(e) => (e.currentTarget.style.borderColor = "var(--border)")}
      />
    </div>
  );
}