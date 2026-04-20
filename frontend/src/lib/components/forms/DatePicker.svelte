<script lang="ts">
  import { DatePicker as DatePickerPrimitive } from "bits-ui";
  import { CalendarIcon, ChevronLeft, ChevronRight } from "lucide-svelte";
  import { 
    DateFormatter, 
    type DateValue, 
    getLocalTimeZone,
    today
  } from "@internationalized/date";
  import { cn } from "$lib/utils/styles";

  interface Props {
    value?: DateValue;
    onValueChange?: (value: DateValue | undefined) => void;
    error?: boolean;
    class?: string;
    placeholder?: string;
    disabled?: boolean;
  }

  let { 
    value = $bindable(), 
    onValueChange,
    error = false,
    class: className,
    placeholder = "Select a date",
    disabled = false
  }: Props = $props();

  const df = new DateFormatter("en-US", {
    dateStyle: "long",
  });
</script>

<DatePickerPrimitive.Root
  bind:value
  onValueChange={onValueChange}
  {disabled}
  weekdayFormat="short"
>
  <div class={cn("grid w-full gap-1.5", className)}>
    <DatePickerPrimitive.Input
      class={cn(
        "flex h-11 w-full items-center rounded-xl border bg-slate-50/50 dark:bg-slate-900/50 px-4 py-2 text-sm ring-offset-background focus-within:outline-none focus-within:ring-2 focus-within:ring-primary/20 focus-within:border-primary disabled:cursor-not-allowed disabled:opacity-50 transition-all duration-200",
        error ? "border-destructive focus-within:ring-destructive/20" : "border-input hover:border-muted-foreground/30"
      )}
    >
      {#snippet children({ segments })}
        {#each segments as { part, value }}
          <DatePickerPrimitive.Segment
            {part}
            class="rounded-sm px-0.5 focus:bg-accent focus:text-accent-foreground outline-none transition-colors"
          >
            {value}
          </DatePickerPrimitive.Segment>
        {/each}
        <DatePickerPrimitive.Trigger
          class="ml-auto inline-flex items-center justify-center rounded-lg p-1 hover:bg-accent hover:text-accent-foreground transition-colors"
        >
          <CalendarIcon class="h-4 w-4 opacity-50" />
        </DatePickerPrimitive.Trigger>
      {/snippet}
    </DatePickerPrimitive.Input>
  </div>

  <DatePickerPrimitive.Content
    sideOffset={6}
    class="z-50 rounded-2xl border bg-popover p-4 text-popover-foreground shadow-2xl animate-in fade-in zoom-in-95 duration-200"
  >
    <DatePickerPrimitive.Calendar>
      {#snippet children({ months, weekdays })}
        <header class="flex items-center justify-between pb-4">
          <DatePickerPrimitive.PrevButton
            class="inline-flex h-8 w-8 items-center justify-center rounded-xl hover:bg-accent transition-colors"
          >
            <ChevronLeft class="h-4 w-4" />
          </DatePickerPrimitive.PrevButton>
          <DatePickerPrimitive.Heading class="text-sm font-semibold" />
          <DatePickerPrimitive.NextButton
            class="inline-flex h-8 w-8 items-center justify-center rounded-xl hover:bg-accent transition-colors"
          >
            <ChevronRight class="h-4 w-4" />
          </DatePickerPrimitive.NextButton>
        </header>

        <div class="flex flex-col space-y-4 pt-4 sm:flex-row sm:space-x-4 sm:space-y-0">
          {#each months as month}
            <DatePickerPrimitive.Grid class="w-full border-collapse space-y-1">
              <DatePickerPrimitive.GridHead>
                <DatePickerPrimitive.GridRow class="flex w-full">
                  {#each weekdays as day}
                    <DatePickerPrimitive.HeadCell
                      class="w-8 rounded-md text-[0.8rem] font-normal text-muted-foreground"
                    >
                      {day.slice(0, 1)}
                    </DatePickerPrimitive.HeadCell>
                  {/each}
                </DatePickerPrimitive.GridRow>
              </DatePickerPrimitive.GridHead>
              <DatePickerPrimitive.GridBody>
                {#each month.weeks as week}
                  <DatePickerPrimitive.GridRow class="flex w-full mt-2">
                    {#each week as date}
                      <DatePickerPrimitive.Cell
                        {date}
                        month={month.value}
                        class="relative p-0 text-center text-sm focus-within:relative focus-within:z-20 [&:has([data-selected])]:bg-accent first:[&:has([data-selected])]:rounded-l-md last:[&:has([data-selected])]:rounded-r-md"
                      >
                        <DatePickerPrimitive.Day
                          class={cn(
                            "inline-flex h-8 w-8 items-center justify-center rounded-xl p-0 font-normal transition-all hover:bg-primary/20 hover:text-primary data-selected:bg-primary data-selected:text-primary-foreground data-selected:font-semibold data-today:bg-accent data-today:text-accent-foreground outline-none"
                          )}
                        />
                      </DatePickerPrimitive.Cell>
                    {/each}
                  </DatePickerPrimitive.GridRow>
                {/each}
              </DatePickerPrimitive.GridBody>
            </DatePickerPrimitive.Grid>
          {/each}
        </div>
      {/snippet}
    </DatePickerPrimitive.Calendar>
  </DatePickerPrimitive.Content>
</DatePickerPrimitive.Root>
