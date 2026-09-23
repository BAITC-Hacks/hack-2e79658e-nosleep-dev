"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { Globe2, Menu, X } from "lucide-react";

type SiteHeaderProps = {
  page: "landing" | "play";
};

export function SiteHeader({ page }: SiteHeaderProps) {
  const [menuOpen, setMenuOpen] = useState(false);
  const [scrolled, setScrolled] = useState(false);
  const isLanding = page === "landing";

  useEffect(() => {
    const updateHeader = () => setScrolled(window.scrollY > 40);
    updateHeader();
    window.addEventListener("scroll", updateHeader, { passive: true });
    return () => window.removeEventListener("scroll", updateHeader);
  }, []);

  const closeMenu = () => setMenuOpen(false);

  return (
    <nav className={`armeta-nav armeta-site-nav ${scrolled ? "is-scrolled" : ""}`} aria-label="Основная навигация">
      <Link className="armeta-logo" href={isLanding ? "#top" : "/"} aria-label="BesSheshim — на главную" onClick={closeMenu}>
        <span className="armeta-logo-mark" aria-hidden="true">BS/5</span>
        <span>BES SHESHIM</span>
      </Link>
      <button className="armeta-menu-button" type="button" aria-expanded={menuOpen} aria-label={menuOpen ? "Закрыть меню" : "Открыть меню"} onClick={() => setMenuOpen((value) => !value)}>
        {menuOpen ? <X size={22} /> : <Menu size={22} />}
      </button>
      <div className={`armeta-nav-links ${menuOpen ? "is-open" : ""}`}>
        <a className="is-active" href={isLanding ? "#tools" : "#simulator"} onClick={closeMenu}>Симулятор</a>
        <Link href={isLanding ? "#about" : "/#about"} onClick={closeMenu}>Как это работает</Link>
        <Link href={isLanding ? "#security" : "/#security"} onClick={closeMenu}>Принципы</Link>
        <Link href={isLanding ? "#footer" : "/#footer"} onClick={closeMenu}>О проекте</Link>
      </div>
      <div className="armeta-nav-actions">
        <span className="armeta-language"><Globe2 size={18} /> рус</span>
        <Link className="armeta-contact" href={isLanding ? "/play" : "/"}>{isLanding ? "Начать игру" : "На главную"}</Link>
      </div>
    </nav>
  );
}
