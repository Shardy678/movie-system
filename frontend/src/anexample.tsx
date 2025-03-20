"use client";

import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

type Seat = {
  id: string;
  row: string;
  column: number;
  isAvailable: boolean;
  isSelected: boolean;
};

type SeatSelectionProps = {
  availableSeats: string[];
};

export function SeatSelection({ availableSeats }: SeatSelectionProps) {
  const [seats, setSeats] = useState<Seat[]>([]);
  const [selectedSeats, setSelectedSeats] = useState<Seat[]>([]);

  useEffect(() => {
    // Parse available seats to determine theater dimensions
    const parsedSeats = availableSeats.map((seatId) => {
      const row = seatId.charAt(0);
      const column = Number.parseInt(seatId.substring(1));
      return {
        id: seatId,
        row,
        column,
        isAvailable: true,
        isSelected: false,
      };
    });

    // Find all unique rows and columns
    const rows = Array.from(
      new Set(parsedSeats.map((seat) => seat.row))
    ).sort();
    const maxColumn = Math.max(...parsedSeats.map((seat) => seat.column));

    // Create complete theater layout including unavailable seats
    const allSeats: Seat[] = [];

    rows.forEach((row) => {
      for (let col = 1; col <= maxColumn; col++) {
        const seatId = `${row}${col}`;
        const existingSeat = parsedSeats.find((seat) => seat.id === seatId);

        if (existingSeat) {
          allSeats.push(existingSeat);
        } else {
          allSeats.push({
            id: seatId,
            row,
            column: col,
            isAvailable: false,
            isSelected: false,
          });
        }
      }
    });

    setSeats(allSeats);
  }, [availableSeats]);

  const toggleSeatSelection = (seatId: string) => {
    setSeats((prevSeats) =>
      prevSeats.map((seat) => {
        if (seat.id === seatId && seat.isAvailable) {
          const newIsSelected = !seat.isSelected;

          // Update selected seats list
          if (newIsSelected) {
            setSelectedSeats((prev) => [...prev, seat]);
          } else {
            setSelectedSeats((prev) => prev.filter((s) => s.id !== seatId));
          }

          return { ...seat, isSelected: newIsSelected };
        }
        return seat;
      })
    );
  };

  // Group seats by row for display
  const seatsByRow = seats.reduce((acc, seat) => {
    if (!acc[seat.row]) {
      acc[seat.row] = [];
    }
    acc[seat.row].push(seat);
    return acc;
  }, {} as Record<string, Seat[]>);

  const rowLabels = Object.keys(seatsByRow).sort();

  return (
    <Card className="w-full max-w-3xl mx-auto">
      <CardHeader>
        <CardTitle className="text-center">Select Your Seats</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="mb-8 flex justify-center">
          <div className="w-full max-w-md">
            {/* Screen */}
            <div className="relative mb-10 mx-auto">
              <div className="h-2 bg-primary rounded-md w-4/5 mx-auto mb-1"></div>
              <p className="text-center text-sm text-muted-foreground">
                Screen
              </p>
            </div>

            {/* Seat map */}
            <div className="space-y-3">
              {rowLabels.map((row) => (
                <div key={row} className="flex items-center">
                  <div className="w-6 text-center font-medium">{row}</div>
                  <div className="flex flex-1 justify-center gap-2">
                    {seatsByRow[row]
                      .sort((a, b) => a.column - b.column)
                      .map((seat) => (
                        <button
                          key={seat.id}
                          onClick={() => toggleSeatSelection(seat.id)}
                          disabled={!seat.isAvailable}
                          className={`
                          w-7 h-7 rounded-t-md text-xs flex items-center justify-center transition-colors
                          ${
                            !seat.isAvailable
                              ? "bg-gray-200 text-gray-400 cursor-not-allowed dark:bg-gray-700"
                              : seat.isSelected
                              ? "bg-primary text-primary-foreground"
                              : "bg-green-100 hover:bg-green-200 dark:bg-green-900/30 dark:hover:bg-green-900/50"
                          }
                        `}
                          aria-label={`Seat ${seat.id}`}
                        >
                          {seat.column}
                        </button>
                      ))}
                  </div>
                </div>
              ))}
            </div>

            {/* Legend */}
            <div className="mt-8 flex justify-center gap-6">
              <div className="flex items-center gap-2">
                <div className="w-5 h-5 rounded-t-md bg-green-100 dark:bg-green-900/30"></div>
                <span className="text-sm">Available</span>
              </div>
              <div className="flex items-center gap-2">
                <div className="w-5 h-5 rounded-t-md bg-primary"></div>
                <span className="text-sm">Selected</span>
              </div>
              <div className="flex items-center gap-2">
                <div className="w-5 h-5 rounded-t-md bg-gray-200 dark:bg-gray-700"></div>
                <span className="text-sm">Unavailable</span>
              </div>
            </div>
          </div>
        </div>
      </CardContent>
      <CardFooter className="flex flex-col gap-2">
        <div className="text-center">
          {selectedSeats.length > 0 ? (
            <p>
              Selected seats: {selectedSeats.map((seat) => seat.id).join(", ")}
            </p>
          ) : (
            <p>No seats selected</p>
          )}
        </div>
        <Button className="w-full" disabled={selectedSeats.length === 0}>
          Continue to Checkout
        </Button>
      </CardFooter>
    </Card>
  );
}
