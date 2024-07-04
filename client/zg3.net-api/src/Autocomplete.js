import React, { useState, useEffect } from 'react';
import axios from 'axios';

const Autocomplete = () => {
  const [query, setQuery] = useState('');
  const [suggestions, setSuggestions] = useState([]);
  const [selectedRestaurant, setSelectedRestaurant] = useState(null);

  useEffect(() => {
    if (query.length > 1) {
      axios.get(`http://localhost:8012/restaurants/autocomplete/${query}`)
        .then(response => {
          setSuggestions(response.data);
        })
        .catch(error => {
          console.error('Error fetching autocomplete suggestions:', error);
        });
    } else {
      setSuggestions([]);
    }
  }, [query]);

  const handleSelect = (restaurant) => {
    setSelectedRestaurant(restaurant);
    setQuery(restaurant[0]);
    setSuggestions([]);
  };

  return (
    <div>
      <input
        type="text"
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        placeholder="Type to search..."
      />
      <ul>
        {suggestions.map((restaurant) => (
          <li key={restaurant[1]} onClick={() => handleSelect(restaurant)}>
            {restaurant[0]}
          </li>
        ))}
      </ul>
      {selectedRestaurant && (
        <div>
          <p>Selected Restaurant:</p>
          <p>Name: {selectedRestaurant[0]}</p>
          <p>ID: {selectedRestaurant[1]}</p>
        </div>
      )}
    </div>
  );
};

export default Autocomplete;
