import { useEffect, useState } from "react";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "./components/ui/card";
import { Button } from "./components/ui/button";
import { toast } from "sonner";

type Seat = {
  id: string;
  row: string;
  column: number;
  isAvailable: boolean;
  isSelected: boolean;
};

function SeatSelection({
  availableSeats,
  movieId,
  showtimeId,
  onClose
}: {
  availableSeats: string[];
  movieId: number;
  showtimeId: number;
  onClose: () => void
}) {
  const [seats, setSeats] = useState<Seat[]>([]);

  useEffect(() => {
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

    const rows = Array.from(
      new Set(parsedSeats.map((seat) => seat.row))
    ).sort();
    const maxColumn = Math.max(...parsedSeats.map((seat) => seat.column));

    const allSeats: Seat[] = [];

    rows.forEach((row) => {
      for (let col = 1; col <= maxColumn; col++) {
        const seatId = `${row}${col}`;
        const existingSeat = parsedSeats.find((seat) => seat.id === seatId);

        allSeats.push(
          existingSeat ?? {
            id: seatId,
            row,
            column: col,
            isAvailable: false,
            isSelected: false,
          }
        );
      }
    });
    setSeats(allSeats);
  }, [availableSeats]);

  const toggleSeatSelection = (seatId: string) => {
    setSeats((prevSeats) =>
      prevSeats.map((seat) =>
        seat.id === seatId && seat.isAvailable
          ? { ...seat, isSelected: !seat.isSelected }
          : seat
      )
    );
  };

  const handleReserveSeats = async () => {
    const selectedSeats = seats
      .filter((seat) => seat.isSelected)
      .map((seat) => seat.id);

    try {
      const token = localStorage.getItem("token");
      const response = await fetch("http://localhost:8080/reserve/add", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          movie_id: movieId,
          showtime_id: showtimeId,
          seats: selectedSeats,
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error("Failed to reserve seats");
      }

      setSeats((prevSeats) =>
        prevSeats.map((seat) =>
          selectedSeats.includes(seat.id)
            ? { ...seat, isAvailable: false, isSelected: false }
            : seat
        )
      );

      toast("Seat(s) reserved");
      onClose()
    } catch (error) {
      console.error("Error reserving seats:", error);
      toast("Failed to reserve seats");
    }
  };

  const selectedSeats = seats.filter((seat) => seat.isSelected);
  const seatsByRow = seats.reduce((acc, seat) => {
    if (!acc[seat.row]) acc[seat.row] = [];
    acc[seat.row].push(seat);
    return acc;
  }, {} as Record<string, Seat[]>);

  return (
    <Card className="w-full max-w-3xl mx-auto">
      <CardHeader>
        <CardTitle className="text-center">Select Your Seats</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="mb-8 flex justify-center">
          <div className="w-full max-w-md">
            <div className="relative mb-10 mx-auto">
              <div className="h-2 bg-primary rounded-md w-4/5 mx-auto mb-1"></div>
              <p className="text-center text-sm text-muted-foreground">
                Screen
              </p>
            </div>

            <div className="space-y-3">
              {Object.keys(seatsByRow)
                .sort()
                .map((row) => (
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
        <Button
          className="w-full"
          disabled={selectedSeats.length === 0}
          onClick={handleReserveSeats}
        >
          Reserve seats
        </Button>
      </CardFooter>
    </Card>
  );
}

export default SeatSelection;
