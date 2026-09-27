import ErrorMessage from "@/components/ErrorMessage"
import { Button } from "@/components/ui/button"
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"
import { useApiClient, type ApiClient } from "@/hooks/use-api-client"
import type { PaginationResult } from "@/lib/dto"
import { cn } from "@/lib/utils"
import { useQuery, type UseQueryOptions } from "@tanstack/react-query"
import { Check, ChevronsUpDown, CircleX, Loader, Loader2 } from "lucide-react"
import { useState } from "react"
import { useTranslation } from "react-i18next"
import { useDebounceCallback } from "usehooks-ts"

interface props<TDto extends { id: string }> {
  id: string
  name: string
  onSelect: (id: string, name: string) => void
  // each resource's queryKey tuple shape differs; only used positionally here
  /* eslint-disable @typescript-eslint/no-explicit-any */
  queries: {
    get: (apiClient: ApiClient, id: string) => UseQueryOptions<TDto, Error, TDto, any>
    list: (
      apiClient: ApiClient,
      searchText: string,
      pageIndex?: number,
      recordsPerPage?: number
    ) => UseQueryOptions<PaginationResult<TDto>, Error, PaginationResult<TDto>, any>
  }
  /* eslint-enable @typescript-eslint/no-explicit-any */
  getLabel: (item: TDto) => string
  placeholder: string
}

// Generic searchable combobox for picking a related record of any table
// exposing `get(apiClient, id)` and `list(apiClient, searchText, pageIndex,
// recordsPerPage)` queryOptions factories (see e.g. project/user/_queries.ts).
const ResourceCombobox = <TDto extends { id: string }>({
  id,
  name,
  onSelect,
  queries,
  getLabel,
  placeholder,
}: props<TDto>) => {
  const { t } = useTranslation()
  const [open, setOpen] = useState(false)
  const apiClient = useApiClient()
  const [searchValue, setSearchValue] = useState("")
  const [debouncedSearchValue, setDebouncedSearchValue] = useState("")

  const debounceSearch = useDebounceCallback((str: string) => {
    setDebouncedSearchValue(str)
  }, 250)

  const {
    data: internalData,
    isLoading: internalDataIsLoading,
    isError: internalDataIsError,
  } = useQuery({
    ...queries.get(apiClient, id),
    enabled: id !== "" && name === "",
  }) // we need to load only if name is not given.

  // ------------------------- Query: List slim -----------------------------------
  const { data, isLoading, isError } = useQuery(
    queries.list(apiClient, debouncedSearchValue, 0, 50) // this is dropdown combo box, so by default list only 50 items
  )

  const searchInputHandler = (enteredText: string) => {
    setSearchValue(enteredText)
    debounceSearch(enteredText)
  }

  const onSelectComponent = (selectedId: string) => {
    if (selectedId === id) {
      onSelect("", "")
    } else {
      const selectedObj = data?.items.find((p) => p.id === selectedId)
      if (selectedObj) {
        onSelect(selectedObj.id, getLabel(selectedObj))
      }
    }
    setOpen(false)
  }

  // ------------------------- Render -----------------------------------
  const renderSelectedItem = () => {
    if (name) {
      return <>{name}</>
    }

    if (internalDataIsLoading) {
      return <Loader className="size-4 animate-spin" />
    }

    if (internalDataIsError) {
      return <CircleX className="size-4 text-red-400" />
    }

    if (internalData) {
      return <>{getLabel(internalData)}</>
    }

    return <>{placeholder}</>
  }

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger
        render={
          <Button
            variant="outline"
            role="combobox"
            aria-expanded={open}
            className="w-full justify-between"
          />
        }
      >
        {isLoading ? (
          <Loader2 className="size-2 animate-spin" />
        ) : (
          renderSelectedItem()
        )}
        <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
      </PopoverTrigger>
      <PopoverContent className="w-full p-0">
        {isError ? (
          <ErrorMessage title="Error!">{t("err-loading-data")}</ErrorMessage>
        ) : (
          <Command shouldFilter={false}>
            <CommandInput
              value={searchValue}
              placeholder={placeholder}
              onValueChange={searchInputHandler}
            />
            <CommandEmpty>{t("no-data")}</CommandEmpty>
            <CommandGroup>
              <CommandList>
                {data?.items.map((row) => (
                  <CommandItem
                    key={row.id}
                    value={row.id}
                    onSelect={onSelectComponent}
                  >
                    <Check
                      className={cn(
                        "mr-2 size-4",
                        id === row.id ? "opacity-100" : "opacity-0"
                      )}
                    />
                    {getLabel(row)}
                  </CommandItem>
                ))}
              </CommandList>
            </CommandGroup>
          </Command>
        )}
      </PopoverContent>
    </Popover>
  )
}

export default ResourceCombobox
