import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import styles from './NewMovie.module.css';

function NewMovie() {
    const [title, setTitle] = useState('');
    const [description, setDescription] = useState('');
    const [genre, setGenre] = useState('');
    const [posterImage, setPosterImage] = useState('');
    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            const token = localStorage.getItem('token');
            if (!token) {
                throw new Error('No authentication token found');
            }

            const response = await fetch('http://localhost:8080/movies/add', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    title,
                    description,
                    genre,
                    poster_image: posterImage
                })
            });

            if (!response.ok) {
                throw new Error('Failed to add movie');
            }

            navigate('/');
        } catch (error) {
            console.error('Error adding movie:', error);
        }
    };

    return (
        <div className={styles.container}>
            <h1 className={styles.title}>Add New Movie</h1>
            <form onSubmit={handleSubmit} className={styles.form}>
                <div className={styles.formGroup}>
                    <label className={styles.label}>Title</label>
                    <input
                        className={styles.input}
                        type="text"
                        value={title}
                        onChange={(e) => setTitle(e.target.value)}
                        required
                    />
                </div>
                <div className={styles.formGroup}>
                    <label className={styles.label}>Description</label>
                    <textarea
                        className={styles.textarea}
                        value={description}
                        onChange={(e) => setDescription(e.target.value)}
                        required
                    />
                </div>
                <div className={styles.formGroup}>
                    <label className={styles.label}>Genre</label>
                    <input
                        className={styles.input}
                        type="text"
                        value={genre}
                        onChange={(e) => setGenre(e.target.value)}
                        required
                    />
                </div>
                <div className={styles.formGroup}>
                    <label className={styles.label}>Poster Image URL</label>
                    <input
                        className={styles.input}
                        type="url"
                        value={posterImage}
                        onChange={(e) => setPosterImage(e.target.value)}
                        required
                    />
                </div>
                <button type="submit" className={styles.button}>
                    Add Movie
                </button>
            </form>
        </div>
    );
}

export default NewMovie;