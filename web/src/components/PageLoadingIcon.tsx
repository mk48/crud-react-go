import { LoaderPinwheel } from "lucide-react";

const PageLoadingIcon: React.FC = () => {
  return (
    <div className="flex min-h-40 items-center justify-center">
      <LoaderPinwheel className="animate-spin" />
    </div>
  );
};

export default PageLoadingIcon;
