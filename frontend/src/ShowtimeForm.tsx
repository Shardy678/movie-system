"use client";

import type React from "react";

import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { format } from "date-fns";

import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { AlertCircle, CalendarIcon } from "lucide-react";
import { cn } from "@/lib/utils";

interface ShowtimeFormProps {
  movieId: number;
  onClose: () => void;
  onShowtimeAdd: (newShowtime: any) => void;
}

function ShowtimeForm({ movieId, onClose, onShowtimeAdd }: ShowtimeFormProps) {
  const [date, setDate] = useState<Date | undefined>(undefined);
  const [time, setTime] = useState("");
  const [capacity, setCapacity] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [validationErrors, setValidationErrors] = useState({
    date: "",
    time: "",
    capacity: "",
  });
  const navigate = useNavigate();

  const validateForm = () => {
    const errors = {
      date: "",
      time: "",
      capacity: "",
    };
    let isValid = true;

    if (!date) {
      errors.date = "Date is required";
      isValid = false;
    }

    if (!time) {
      errors.time = "Time is required";
      isValid = false;
    }

    if (!capacity || capacity <= 0) {
      errors.capacity = "Capacity must be at least 1";
      isValid = false;
    }

    setValidationErrors(errors);
    return isValid;
  };

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    const dateTime = new Date(date!);
    const [hours, minutes] = time.split(":").map(Number);
    dateTime.setHours(hours, minutes);

    const showtimeData = {
      movie_id: movieId,
      start_time: dateTime.toISOString(),
      capacity,
    };

    try {
      const token = localStorage.getItem("token");
      if (!token) {
        navigate("/login");
        return;
      }
      // Не возвращает айди корректно, поэтому после создания в стейте у нового шоутайма айди: 0
      const response = await fetch("http://localhost:8080/showtimes/add", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(showtimeData),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || "Failed to add showtime.");
      }

      const newShowtime = await response.json();
      onShowtimeAdd(newShowtime);
      onClose();

      toast("Showtime added successfully", {
        description: `Showtime for movie ID ${movieId} has been added.`,
      });
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred.");
    }
  }

  return (
    <>
      <Dialog open={true} onOpenChange={(open) => !open && onClose()}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Add Showtime</DialogTitle>
          </DialogHeader>

          {error && (
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="date">Date</Label>
              <Popover>
                <PopoverTrigger asChild>
                  <Button
                    id="date"
                    variant="outline"
                    className={cn(
                      "w-full justify-start text-left font-normal",
                      !date && "text-muted-foreground"
                    )}
                  >
                    <CalendarIcon className="mr-2 h-4 w-4" />
                    {date ? format(date, "PPP") : <span>Select date</span>}
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="w-auto p-0" align="start">
                  <Calendar
                    mode="single"
                    selected={date}
                    onSelect={setDate}
                    initialFocus
                  />
                </PopoverContent>
              </Popover>
              {validationErrors.date && (
                <p className="text-sm font-medium text-destructive">
                  {validationErrors.date}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="time">Time</Label>
              <Input
                id="time"
                type="time"
                value={time}
                onChange={(e) => setTime(e.target.value)}
                required
              />
              {validationErrors.time && (
                <p className="text-sm font-medium text-destructive">
                  {validationErrors.time}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="capacity">Capacity</Label>
              <Input
                id="capacity"
                type="number"
                placeholder="Capacity"
                value={capacity || ""}
                onChange={(e) => setCapacity(Number(e.target.value))}
                required
              />
              {validationErrors.capacity && (
                <p className="text-sm font-medium text-destructive">
                  {validationErrors.capacity}
                </p>
              )}
            </div>

            <DialogFooter className="gap-2 sm:gap-0">
              <Button type="button" variant="outline" onClick={onClose}>
                Cancel
              </Button>
              <Button type="submit">Submit</Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </>
  );
}

export default ShowtimeForm;
