import { useState } from 'react';
import styles from './ShowtimeForm.module.css';

interface ShowtimeFormProps {
  movieId: number;
  onClose: () => void;
  onAddShowtime: (showtime: { movie_id: number; start_time: string; capacity: number }) => void;
}

function ShowtimeForm({ movieId, onClose, onAddShowtime }: ShowtimeFormProps) {
  const [startTime, setStartTime] = useState('');
  const [capacity, setCapacity] = useState(0);

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();

    if (!startTime || capacity <= 0) {
      alert('Please enter a valid showtime and capacity.');
      return;
    }

    onAddShowtime({ movie_id: movieId, start_time: startTime, capacity });
    onClose(); 
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
          value={capacity || ''}
          onChange={(e) => setCapacity(Number(e.target.value))}
          required
        />
        <div className={styles.buttonGroup}>
          <button type="submit">Submit</button>
          <button type="button" onClick={onClose}>Cancel</button>
        </div>
      </form>
    </div>
  );
}

export default ShowtimeForm;
