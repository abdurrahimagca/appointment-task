import * as React from "react"

function Card({ className = "", ...props }: React.ComponentProps<"div">) {
  return <div className={className} {...props} />
}

function CardHeader({ className = "", ...props }: React.ComponentProps<"div">) {
  return <div className={className} {...props} />
}

function CardTitle({ className = "", ...props }: React.ComponentProps<"div">) {
  return <div className={className} {...props} />
}

function CardDescription({ className = "", ...props }: React.ComponentProps<"div">) {
  return <div className={className} {...props} />
}

function CardAction({ className = "", ...props }: React.ComponentProps<"div">) {
  return <div className={className} {...props} />
}

function CardContent({ className = "", ...props }: React.ComponentProps<"div">) {
  return <div className={className} {...props} />
}

function CardFooter({ className = "", ...props }: React.ComponentProps<"div">) {
  return <div className={className} {...props} />
}

export {
  Card,
  CardHeader,
  CardFooter,
  CardTitle,
  CardAction,
  CardDescription,
  CardContent,
}
