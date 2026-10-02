interface props {
  values: Record<string, unknown>
}

// A compact key/value list for an operation's free-form JSON (metadata,
// client info).
const KeyValueList: React.FC<props> = ({ values }) => (
  <dl className="grid grid-cols-[auto_1fr] gap-x-4">
    {Object.entries(values).map(([key, value]) => (
      <div key={key} className="contents">
        <dt className="text-muted-foreground">{key}</dt>
        <dd className="break-all">
          {typeof value === "string" ? value : JSON.stringify(value)}
        </dd>
      </div>
    ))}
  </dl>
)

export default KeyValueList
