import { useState } from "react";
import styles from "./ShowtimeForm.module.css";
import { useNavigate } from "react-router-dom";

interface ShowtimeFormProps {
  movieId: number;
  onClose: () => void;
  onAddShowtime: (showtime: {
    movie_id: number;
    start_time: string;
    capacity: number;
  }) => void;
}

function ShowtimeForm({ movieId, onClose, onAddShowtime }: ShowtimeFormProps) {
  const [startTime, setStartTime] = useState("");
  const [capacity, setCapacity] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();
  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();

    if (!startTime || capacity <= 0) {
      alert("Please enter a valid showtime and capacity.");
      return;
    }

    const showtimeData = {
      movie_id: movieId,
      start_time: new Date(startTime).toISOString(),
      capacity,
    };

    try {
      const token = localStorage.getItem("token");
      if (!token) {
        navigate("/login");
        return;
      }

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
      onAddShowtime(newShowtime);
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred.");
    }
  }

  if (error) {
    return <div>Error</div>;
  }

  return (
    <div className={styles.formOverlay}>
      <form className={styles.showtimeForm} onSubmit={handleSubmit}>
        <h3>Add Showtime</h3>
        <label>Showtime Date & Time:</label>
        <input
          type="datetime-local"
          value={startTime}
          onChange={(e) => setStartTime(e.target.value)}
          required
        />
        <label>Capacity:</label>
        <input
          type="number"
          placeholder="Capacity"
          value={capacity || ""}
          onChange={(e) => setCapacity(Number(e.target.value))}
          required
        />
        <div className={styles.buttonGroup}>
          <button type="submit">Submit</button>
          <button type="button" onClick={onClose}>
            Cancel
          </button>
        </div>
      </form>
    </div>
  );
}

export default ShowtimeForm;
