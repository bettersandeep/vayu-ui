import * as Sentry from "@sentry/react"
import { StrictMode } from "react"
import { createRoot } from "react-dom/client"

import "./index.css"
import App from "@/app/App.tsx"

// Async Sentry init from server-provided runtime config; never blocks first render.
fetch("/api/sentry")
	.then(res => (res.ok ? res.json() : null))
	.then(cfg => {
		if (cfg?.dsn) {
			Sentry.init({
				dsn: cfg.dsn,
				release: cfg.release || undefined,
				tracesSampleRate: 0,
			})
		}
	})
	.catch(() => {
		// Sentry stays disabled if config fetch fails
	})

createRoot(document.getElementById("root")!).render(
	<StrictMode>
		<Sentry.ErrorBoundary
			fallback={
				<div className="flex h-screen items-center justify-center">
					Something went wrong
				</div>
			}
		>
			<App />
		</Sentry.ErrorBoundary>
	</StrictMode>,
)
