import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Check } from "lucide-react";

interface props {
  children: React.ReactNode;
  title: string;
}

const SuccessMessage: React.FC<props> = ({ title, children }) => {
  return (
    <Alert className="mb-2 mt-2 bg-green-200">
      <Check className="size-4" />
      <AlertTitle>{title}</AlertTitle>
      <AlertDescription> {children}</AlertDescription>
    </Alert>
  );
};

export default SuccessMessage;
