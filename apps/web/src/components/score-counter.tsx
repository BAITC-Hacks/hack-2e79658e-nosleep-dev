"use client";

import { useEffect, useRef, useState } from "react";

type ScoreCounterProps = {
  value: number;
  decimals?: number;
  duration?: number;
  className?: string;
};

export function ScoreCounter({ value, decimals = 2, duration = 1100, className }: ScoreCounterProps) {
  const elementRef = useRef<HTMLSpanElement>(null);
  const [displayValue, setDisplayValue] = useState(value);

  useEffect(() => {
    const element = elementRef.current;
    if (!element) return;

    const motion = window.matchMedia("(prefers-reduced-motion: reduce)");
    if (motion.matches) {
      return;
    }

    let frame = 0;
    let startedAt = 0;
    let hasRun = false;

    const run = () => {
      if (hasRun) return;
      hasRun = true;
      setDisplayValue(0);

      const tick = (now: number) => {
        if (!startedAt) startedAt = now;
        const progress = Math.min((now - startedAt) / duration, 1);
        const eased = 1 - Math.pow(1 - progress, 3);
        setDisplayValue(value * eased);
        if (progress < 1) frame = window.requestAnimationFrame(tick);
      };

      frame = window.requestAnimationFrame(tick);
    };

    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          run();
          observer.disconnect();
        }
      },
      { threshold: 0.45 },
    );

    observer.observe(element);
    return () => {
      observer.disconnect();
      window.cancelAnimationFrame(frame);
    };
  }, [duration, value]);

  return (
    <span ref={elementRef} className={className} aria-label={`${value.toFixed(decimals)} points`}>
      {displayValue.toFixed(decimals)}
    </span>
  );
}
