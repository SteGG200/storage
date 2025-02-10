import { ArrowRight } from "lucide-react";
import Link from "next/link";

export default function Home() {
  return (
    <main className="min-h-dvh flex flex-col justify-center items-center space-y-8">
      <p className="text-5xl">Storage Management App</p>
      <Link href="/storage" className="flex items-center">
        <p className="text-xl after:content-space hover:after:content-none">Click here to enter the storage</p>
        <ArrowRight />
      </Link>
    </main>
  )
}
