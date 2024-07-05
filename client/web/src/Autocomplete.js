import React, { useState, useEffect, useRef } from 'react';
import axios from 'axios';
import './Autocomplete.css'; // Create and import a CSS file for styling

const Autocomplete = () => {
  const [query, setQuery] = useState('');
  const [suggestions, setSuggestions] = useState([]);
  const [selectedRestaurant, setSelectedRestaurant] = useState(null);
  const [isSelecting, setIsSelecting] = useState(false);
  const dropdownRef = useRef(null);

  useEffect(() => {
    if (!isSelecting && query.length > 1) {
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
  }, [query, isSelecting]);

  const handleSelect = (restaurant) => {
    setSelectedRestaurant(restaurant);
    setQuery(restaurant[0]);
    setSuggestions([]);
    setIsSelecting(true);
  };

  const handleInputChange = (e) => {
    setQuery(e.target.value);
    setIsSelecting(false);
  };

  return (
    <div className="autocomplete-container">
      <input
        type="text"
        value={query}
        onChange={handleInputChange}
        placeholder="Type to search..."
        className="autocomplete-input"
        ref={dropdownRef}
      />
      {suggestions.length > 0 && (
        <ul className="autocomplete-suggestions" style={{ top: dropdownRef.current?.offsetHeight }}>
          {suggestions.map((restaurant) => (
            <li key={restaurant[1]} onClick={() => handleSelect(restaurant)} className="autocomplete-suggestion">
              {restaurant[0]}
            </li>
          ))}
        </ul>
      )}
      {selectedRestaurant && (
        <div className="restaurant-details">
          <h2>Restaurant Details</h2>
          <p>Name: {selectedRestaurant[0]}</p>
          <p>ID: {selectedRestaurant[1]}</p>
        </div>
      )}
    </div>
  );
};

export default Autocomplete;
