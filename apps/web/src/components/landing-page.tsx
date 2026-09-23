"use client";

import { useEffect, useState } from "react";
import Image from "next/image";
import Link from "next/link";
import {
  ArrowUpRight, ChartNoAxesCombined, Check,
  LockKeyhole, ShieldCheck,
} from "lucide-react";

import { ArrowFillButton } from "@/components/block/arrow-fill-button";
import { SiteHeader } from "@/components/site-header";
import { api } from "@/lib/api-client";

type Product = {
  title: string;
  body: string;
  image: string;
  accent?: boolean;
};

const products: Product[] = [
  { title: "Пять решений", body: "Распределите единый бюджет между важными городскими инициативами.", image: "/brand/tools/budget.png", accent: true },
  { title: "Город отвечает", body: "Сразу увидьте, как каждый выбор влияет на районы и качество жизни.", image: "/brand/tools/city-impact.png" },
  { title: "Понятные выводы", body: "Получите ясное объяснение сильных сторон, рисков и компромиссов.", image: "/brand/tools/ai-analysis.png" },
  { title: "Живая карта", body: "Следите за изменениями в шести районах Астаны после каждого решения.", image: "/brand/tools/district-map.png" },
  { title: "Честный результат", body: "Все оценки рассчитывает открытая и воспроизводимая модель.", image: "/brand/tools/transparent-score.png" },
  { title: "Ваш ход", body: "Соберите свой план и сравните его с другими сценариями.", image: "/brand/tools/scenario.png" },
];

export function LandingPage() {
  const [health, setHealth] = useState<"loading" | "ok" | "error">("loading");

  useEffect(() => {
    let active = true;
    void api.health()
      .then((response) => { if (active) setHealth(response.status === "ok" ? "ok" : "error"); })
      .catch(() => { if (active) setHealth("error"); });
    return () => { active = false; };
  }, []);

  return (
    <div className="armeta-landing">
      <SiteHeader page="landing" />

      <header className="armeta-hero" id="top">
        <video
          className="armeta-hero-video"
          autoPlay
          muted
          loop
          playsInline
          preload="auto"
          poster="/brand/hero-emerald-poster.jpg"
          aria-hidden="true"
        >
          <source src="/brand/hero-emerald.mp4" type="video/mp4" />
        </video>

        <div className="armeta-hero-copy">
          <p>AI CITY BUDGET SIMULATOR</p>
          <h1><span>Решения, из которых</span><span>строится город.</span></h1>
          <span>Управляйте бюджетом Астаны — прозрачно, наглядно и с реальными последствиями</span>
          <ArrowFillButton
            href="/play"
            animated={false}
            fillBgColor="#70d9b2"
            fillTextColor="#082018"
            hoverFillBgColor="#063d31"
            hoverFillTextColor="#ffffff"
            hoverArrowColor="#ffffff"
          >
            Распределить бюджет
          </ArrowFillButton>
        </div>

        <div className="armeta-proof-row" aria-label="Ключевые факты">
          <div className="armeta-partner-card">
            <div className="armeta-seal">BS</div>
            <p><strong>Пять часов у руля</strong><span>Проверьте свои решения</span></p>
          </div>
          <div className="armeta-stat-card">
            <div><strong>100</strong><span>единиц общего бюджета</span></div>
            <div><strong>05</strong><span>решений в каждом сценарии</span></div>
            <div><strong>05</strong><span>районов на живой карте</span></div>
          </div>
        </div>
      </header>

      <main className="armeta-main">
        <section className="armeta-tools armeta-grid-section" id="tools" aria-labelledby="tools-title">
          <div className="armeta-section-heading">
            <span>ПРОДУКТОВАЯ ЛИНЕЙКА</span>
            <h2 id="tools-title">Умные инструменты для<br />понятных городских решений</h2>
          </div>
          <div className="armeta-product-grid">
            {products.map((product) => (
              <Link className={`armeta-product-card${product.accent ? " is-accent" : ""}`} href="/play" key={product.title}>
                <div className="armeta-product-copy">
                  <h3>{product.title}</h3>
                  <p>{product.body}</p>
                  <span className="armeta-card-link">Попробовать бесплатно <ArrowUpRight size={16} /></span>
                </div>
                <Image className="armeta-product-art" src={product.image} alt="" width={960} height={640} aria-hidden="true" />
              </Link>
            ))}
          </div>
        </section>

        <section className="armeta-about armeta-grid-section" id="about" aria-labelledby="about-title">
          <div className="armeta-section-heading">
            <span>BES SHESHIM</span>
            <h2 id="about-title">Локальный выбор. <em>Общий результат.</em></h2>
            <p>BesSheshim превращает сложную городскую модель в понятный разговор о приоритетах. Каждое решение видно, каждое последствие можно объяснить.</p>
          </div>
          <div className="armeta-country-grid">
            <article><span>ГОРОД</span><h3>Астана как единая система</h3><p>Пять районов связаны общим бюджетом и качеством жизни.</p></article>
            <article><span>ЖИТЕЛИ</span><h3>Решения с человеческим масштабом</h3><p>Сценарий показывает, кто выигрывает и где остаются риски.</p></article>
            <article><span>БУДУЩЕЕ</span><h3>Сравнивайте до того, как выбирать</h3><p>Пробуйте альтернативы и находите более устойчивый баланс.</p></article>
          </div>
          <p className="armeta-about-note">Одна модель для всех сценариев — <strong>единые правила, прозрачные расчёты и сопоставимые результаты.</strong></p>
        </section>

        <section className="armeta-security armeta-grid-section" id="security" aria-labelledby="security-title">
          <div className="armeta-security-intro">
            <div>
              <span className="armeta-pill">ПРИНЦИПЫ</span>
              <h2 id="security-title">Прозрачность модели —<br />наш главный приоритет</h2>
              <p>Числа рассчитываются детерминированным движком. ИИ не меняет результат — он помогает объяснить его простым языком.</p>
              <Link href="/play">Проверить на своём сценарии</Link>
            </div>
            <div className="armeta-trust-lockup"><ShieldCheck size={54} /><span>ПРОВЕРЯЕМО<br /><strong>И ПОНЯТНО</strong></span></div>
          </div>
          <div className="armeta-principles">
            <article><Check size={19} /><h3>Единые правила</h3><p>Каждый сценарий проходит через одну и ту же формулу оценки.</p></article>
            <article><LockKeyhole size={19} /><h3>Без скрытых решений</h3><p>ИИ объясняет расчёт, но не подменяет его своими выводами.</p></article>
            <article><ChartNoAxesCombined size={19} /><h3>Сравнимые результаты</h3><p>Изменения видны на общей шкале и по каждому району отдельно.</p></article>
          </div>
        </section>
      </main>

      <footer className="armeta-footer" id="footer">
        <div className="armeta-footer-top">
          <div><Link className="armeta-logo footer-logo" href="#top"><span className="armeta-logo-mark">BS/5</span><span>BES SHESHIM</span></Link><p>Решения, из которых строится город.<br />AI city budget simulator.</p></div>
          <div><span>ПРОДУКТ</span><Link href="/play">Симулятор</Link><a href="#about">Как это работает</a><a href="#security">Принципы</a></div>
          <div><span>СИСТЕМА</span><p className={`armeta-system ${health}`}><i /> API {health}</p><a href="#tools">Инициативы</a><a href="#about">О модели</a></div>
          <div><span>КОНТАКТ</span><a href="mailto:hello@bessheshim.kz">hello@bessheshim.kz</a><p>Астана, Казахстан</p></div>
        </div>
        <div className="armeta-footer-bottom"><span>© 2026 BesSheshim</span><span>Five decisions. One city.</span></div>
      </footer>
    </div>
  );
}
