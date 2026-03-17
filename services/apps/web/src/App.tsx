import { Route, Routes, useParams } from "react-router"
import { MoonIcon, SunIcon } from "lucide-react"
import { useTheme } from "@/components/theme-provider"
import { Provider } from "./features/appointment/components/provider"

const ThemeToggle = () => {
  const { theme, setTheme } = useTheme()

  const toggle = () => {
    setTheme(theme === "dark" ? "light" : "dark")
  }

  return (
    <button
      aria-label="Toggle theme"
      className="fixed top-4 right-4 z-50 rounded-full border border-border bg-card p-2 text-foreground transition-colors hover:bg-muted"
      onClick={toggle}
      type="button"
    >
      {theme === "dark" ? (
        <SunIcon className="size-5" />
      ) : (
        <MoonIcon className="size-5" />
      )}
    </button>
  )
}

const ProviderRoute = () => {
  const { username = "" } = useParams<{ username: string }>()

  return (
    <main className="min-h-screen bg-background text-foreground">
      <div className="mx-auto flex min-h-screen w-full max-w-7xl flex-col px-4 py-8 sm:px-6 lg:px-8 lg:py-12">
        <Provider username={username} />
      </div>
    </main>
  )
}

const NotFoundRoute = () => {
  return (
    <main className="min-h-screen bg-background text-foreground">
      <div className="mx-auto flex min-h-screen w-full max-w-7xl flex-col px-4 py-8 sm:px-6 lg:px-8 lg:py-12">
        <h1 className="text-2xl font-bold">Not Found</h1>
      </div>
    </main>
  )
}

const App = () => {
  return (
    <>
      <ThemeToggle />
      <Routes>
        <Route path="/:username" element={<ProviderRoute />} />
        <Route path="*" element={<NotFoundRoute />} />
      </Routes>
    </>
  )
}

export default App
