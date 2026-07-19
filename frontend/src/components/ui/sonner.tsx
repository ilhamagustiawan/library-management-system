"use client";

import {
  CircleCheck,
  Info,
  LoaderCircle,
  OctagonX,
  TriangleAlert,
} from "lucide-react";
import { Toaster as Sonner, type ToasterProps } from "sonner";

export function Toaster(props: ToasterProps) {
  return (
    <Sonner
      theme="light"
      className="toaster group"
      icons={{
        success: <CircleCheck />,
        info: <Info />,
        warning: <TriangleAlert />,
        error: <OctagonX />,
        loading: <LoaderCircle className="animate-spin motion-reduce:animate-none" />,
      }}
      {...props}
    />
  );
}
