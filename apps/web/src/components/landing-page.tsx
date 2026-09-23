"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { ArrowDownRight, ArrowRight, MoveUpRight } from "lucide-react";

import { DottedGrid } from "@/components/block/dotted-grid";
import { RectangularTextReveal } from "@/components/block/rectangular-text-reveal";
import { ScoreCounter } from "@/components/score-counter";
import { ScrollReveal } from "@/components/scroll-reveal";
import { api, USE_MOCKS } from "@/lib/api-client";

const EXAMPLE_HREF = "/play?d=M7:nura,M8:nura,M10:nura,M12,M5:saryarka";

const stats = [
  { value: "100", label: "budget units" },
  { value: "05", label: "decisions" },
  { value: "14", label: "initiatives" },
  { value: "05", label: "districts" },
];

const steps = [
  {
    number: "01",
    label: "Allocate",
    title: "Build a five-part plan.",
    body: "Spend one shared 100-unit budget across a catalog of measures. Every player starts with the same city and the same constraint.",
  },
  {
    number: "02",
    label: "Observe",
    title: "Watch the map react.",
    body: "Districts recolor with every pick. Critical indicators below 40 move from red toward green before you submit.",
  },
  {
    number: "03",
    label: "Understand",
    title: "Read the tradeoffs.",
    body: "A deterministic engine owns every number. AI turns the result into a plain-language read of strengths, risks, and compromises.",
  },
];

export function LandingPage() {
  const [health, setHealth] = useState<"loading" | "ok" | "error">("loading");

  useEffect(() => {
    let active = true;
    void api.health()
      .then((response) => {
        if (active) setHealth(response.status === "ok" ? "ok" : "error");
      })
      .catch(() => {
        if (active) setHealth("error");
      });
    return () => { active = false; };
  }, []);

  return (
    <div className="landing-page">
      <nav className="nav-bar" aria-label="Primary navigation">
        <Link className="wordmark" href="#top"><span />MAYOR / 5H</Link>
        <div className="nav-meta">
          <span className={`status-dot ${health}`} /> API {USE_MOCKS ? "mock" : health}
          <a href="#method">How it works</a>
          <Link href="/play">Play</Link>
        </div>
      </nav>

      <main>
        <header className="landing-hero" id="top">
          <DottedGrid className="hero-grid" paused={false} />
          <div className="hero-vignette" aria-hidden="true" />
          <div className="landing-hero-inner">
            <div className="eyebrow hero-eyebrow">CITY BUDGET LAB / ASTANA</div>
            <RectangularTextReveal
              as="h1"
              className="landing-headline"
              baseColor="#f36458"
              overlayColor="#0b0b0b"
              stagger={0.14}
              playOnMount
            >
              Five choices.<br /><span>A city reacts.</span>
            </RectangularTextReveal>
            <div className="landing-hero-bottom">
              <p>Shape one city with one shared budget. See the human cost of every tradeoff before the decision is final.</p>
              <div className="hero-actions">
                <Link className="hero-cta" href="/play">Start a scenario <ArrowRight size={18} /></Link>
                <Link className="text-link" href={EXAMPLE_HREF}>See the worked example <MoveUpRight size={16} /></Link>
              </div>
            </div>
          </div>
          <a className="hero-scroll" href="#proof" aria-label="Scroll to key facts"><span>Scroll to proof</span><ArrowDownRight size={17} /></a>
        </header>

        <ScrollReveal className="stats-strip" delay={60}>
          <section id="proof" aria-label="Scenario facts">
            {stats.map((stat) => (
              <div className="stat" key={stat.label}>
                <strong>{stat.value}</strong>
                <span>{stat.label}</span>
              </div>
            ))}
          </section>
        </ScrollReveal>

        <section className="baseline-section" aria-labelledby="baseline-title">
          <div className="baseline-copy">
            <div className="eyebrow">BASELINE / BEFORE YOUR FIRST MOVE</div>
            <h2 id="baseline-title">The city starts<br /><span>out of balance.</span></h2>
            <p>Nura is the most exposed district, with two indicators below the critical line. The score makes that inequality impossible to hide.</p>
            <Link className="text-link" href={EXAMPLE_HREF}>Follow the 56.54 scenario <MoveUpRight size={16} /></Link>
          </div>
          <ScrollReveal className="baseline-score-card" delay={100}>
            <div className="panel-heading"><span>City today</span><span className="live-label"><i className="pulse" /> Live baseline</span></div>
            <div className="baseline-score">
              <ScoreCounter value={52.56} />
              <small>/ 100</small>
            </div>
            <div className="baseline-card-footer">
              <span>Nura</span>
              <strong><i /> 2 critical indicators</strong>
            </div>
          </ScrollReveal>
        </section>

        <section className="method-section" id="method" aria-labelledby="method-title">
          <div className="method-heading">
            <div>
              <div className="eyebrow">HOW IT WORKS / THREE MOVES</div>
              <h2 id="method-title">Decide. Watch.<br />Understand.</h2>
            </div>
            <p>A short simulation with a visible chain of cause and effect—from budget allocation to district-level impact.</p>
          </div>
          <div className="method-grid">
            {steps.map((step, index) => (
              <ScrollReveal className="method-card" delay={index * 90} key={step.number}>
                <div className="method-card-top"><span>{step.number}</span><span>{step.label}</span></div>
                <h3>{step.title}</h3>
                <p>{step.body}</p>
              </ScrollReveal>
            ))}
          </div>
        </section>

        <section className="closing-cta" aria-labelledby="cta-title">
          <div className="eyebrow">YOUR FIVE HOURS START HERE</div>
          <h2 id="cta-title">What would you<br /><span>change first?</span></h2>
          <div className="closing-row">
            <p>One budget. Five decisions. A city that answers back immediately.</p>
            <Link className="hero-cta closing-button" href="/play">Run the city <ArrowRight size={18} /></Link>
          </div>
        </section>
      </main>
    </div>
  );
}
