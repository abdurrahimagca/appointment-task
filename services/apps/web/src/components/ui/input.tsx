import * as React from "react"

function Input({ className = "", ...props }: React.ComponentProps<"input">) {
  return <input className={className} {...props} />
}

export { Input }
