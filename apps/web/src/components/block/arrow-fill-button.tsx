"use client";

import type { CSSProperties, ComponentProps, ReactNode } from "react";
import Link from "next/link";

import { cn } from "@/lib/utils";
import "./arrow-fill-button.css";

type ArrowFillButtonProps = Omit<ComponentProps<typeof Link>, "children" | "className" | "style"> & {
  children?: ReactNode;
  className?: string;
  style?: CSSProperties;
  bgColor?: string;
  textColor?: string;
  fillBgColor?: string;
  fillTextColor?: string;
  hoverFillBgColor?: string;
  hoverFillTextColor?: string;
  arrowColor?: string;
  hoverArrowColor?: string;
  animated?: boolean;
};

type ArrowFillButtonStyle = CSSProperties & Record<`--btn-${string}`, string>;

export function ArrowFillButton({
  children = "Распределить бюджет",
  className,
  bgColor = "#ffffff",
  textColor = "#111111",
  fillBgColor = "#f36458",
  fillTextColor = "#111111",
  hoverFillBgColor = "#f36458",
  hoverFillTextColor = "#111111",
  arrowColor,
  hoverArrowColor,
  animated = true,
  style,
  ...props
}: ArrowFillButtonProps) {
  const buttonStyle: ArrowFillButtonStyle = {
    "--btn-bg": bgColor,
    "--btn-text": textColor,
    "--btn-fill-bg": fillBgColor,
    "--btn-fill-text": fillTextColor,
    "--btn-fill-bg-hover": hoverFillBgColor,
    "--btn-fill-text-hover": hoverFillTextColor,
    "--btn-arrow": arrowColor ?? fillTextColor,
    "--btn-arrow-hover": hoverArrowColor ?? hoverFillTextColor,
    ...style,
  };

  return (
    <Link
      {...props}
      className={cn("obsidian-arrow-fill-btn", !animated && "is-static", className)}
      style={buttonStyle}
    >
      <span className="obsidian-arrow-fill-btn__text">{children}</span>

      <span aria-hidden="true" className="obsidian-arrow-fill-btn__circle">
        <span>{children}</span>
        <span className="obsidian-arrow-fill-btn__circle-text">
          <svg
            aria-hidden="true"
            className="obsidian-arrow-fill-btn__icon"
            fill="none"
            viewBox="0 0 10 10"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              className="obsidian-arrow-fill-btn__path"
              clipRule="evenodd"
              d="M3.82475e-07 5.625L7.625 5.625L4.125 9.125L5 10L10 5L5 -4.37114e-07L4.125 0.874999L7.625 4.375L4.91753e-07 4.375L3.82475e-07 5.625Z"
              fillRule="evenodd"
            />
            <path
              className="obsidian-arrow-fill-btn__path"
              clipRule="evenodd"
              d="M3.82475e-07 5.625L7.625 5.625L4.125 9.125L5 10L10 5L5 -4.37114e-07L4.125 0.874999L7.625 4.375L4.91753e-07 4.375L3.82475e-07 5.625Z"
              fillRule="evenodd"
            />
          </svg>
        </span>
      </span>
    </Link>
  );
}
