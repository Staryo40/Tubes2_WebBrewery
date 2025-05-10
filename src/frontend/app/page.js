'use client'
import { useState } from "react";
import Navbar from "@/app/components/Navbar";
import Elements from "@/app/components/Elements";
import Visualization from "@/app/components/Visualization";

export default function Home() {
  const [activePage, setActivePage] = useState("elements");

  return (
    <div className="min-h-screen bg-[#260026]">
      <Navbar onNavigate={setActivePage} active={activePage} />
      {activePage === "elements" && <Elements />}
      {activePage === "visualization" && <Visualization />}
    </div>
  );
}
