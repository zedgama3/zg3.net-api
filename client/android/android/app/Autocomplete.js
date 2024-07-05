import React, { useState, useEffect } from 'react';
import { View, TextInput, FlatList, Text, TouchableOpacity, StyleSheet, useColorScheme, TouchableWithoutFeedback, Keyboard } from 'react-native';
import axios from 'axios';

const Autocomplete = () => {
  const [query, setQuery] = useState('');
  const [suggestions, setSuggestions] = useState([]);
  const [selectedRestaurant, setSelectedRestaurant] = useState(null);
  const [isSelecting, setIsSelecting] = useState(false);
  const colorScheme = useColorScheme();

  useEffect(() => {
    if (!isSelecting && query.length > 1) {
      axios.get(`http://10.128.1.5:8012/restaurants/autocomplete/${query}`)
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
    Keyboard.dismiss();
  };

  const handleInputChange = (text) => {
    setQuery(text);
    setIsSelecting(false);
  };

  const renderItem = ({ item }) => (
    <TouchableOpacity onPress={() => handleSelect(item)} style={styles.suggestionItem}>
      <Text style={colorScheme === 'dark' ? styles.suggestionTextDark : styles.suggestionText}>{item[0]}</Text>
    </TouchableOpacity>
  );

  return (
    <TouchableWithoutFeedback onPress={() => Keyboard.dismiss()}>
      <View style={styles.container}>
        <TextInput
          style={colorScheme === 'dark' ? styles.inputDark : styles.input}
          value={query}
          onChangeText={handleInputChange}
          placeholder="Type to search..."
          placeholderTextColor={colorScheme === 'dark' ? '#888' : '#ccc'}
        />
        {suggestions.length > 0 && (
          <FlatList
            data={suggestions}
            renderItem={renderItem}
            keyExtractor={(item) => item[1]}
            style={colorScheme === 'dark' ? styles.suggestionsListDark : styles.suggestionsList}
            keyboardShouldPersistTaps='always'
          />
        )}
        {selectedRestaurant && (
          <View style={styles.details}>
            <Text style={colorScheme === 'dark' ? styles.detailsTextDark : styles.detailsText}>Selected Restaurant:</Text>
            <Text style={colorScheme === 'dark' ? styles.detailsTextDark : styles.detailsText}>Name: {selectedRestaurant[0]}</Text>
            <Text style={colorScheme === 'dark' ? styles.detailsTextDark : styles.detailsText}>ID: {selectedRestaurant[1]}</Text>
          </View>
        )}
      </View>
    </TouchableWithoutFeedback>
  );
};

const styles = StyleSheet.create({
  container: {
    padding: 20,
    flex: 1,
  },
  input: {
    height: 40,
    borderColor: '#ccc',
    borderWidth: 1,
    paddingLeft: 8,
    backgroundColor: '#fff',
    color: '#000',
  },
  inputDark: {
    height: 40,
    borderColor: '#444',
    borderWidth: 1,
    paddingLeft: 8,
    backgroundColor: '#333',
    color: '#fff',
  },
  suggestionsList: {
    borderColor: '#ccc',
    borderWidth: 1,
    marginTop: 5,
    backgroundColor: '#fff',
  },
  suggestionsListDark: {
    borderColor: '#444',
    borderWidth: 1,
    marginTop: 5,
    backgroundColor: '#333',
  },
  suggestionItem: {
    padding: 10,
    borderBottomColor: '#ccc',
    borderBottomWidth: 1,
  },
  suggestionText: {
    color: '#000',
  },
  suggestionTextDark: {
    color: '#fff',
  },
  details: {
    marginTop: 20,
  },
  detailsText: {
    fontSize: 16,
    color: '#000',
  },
  detailsTextDark: {
    fontSize: 16,
    color: '#fff',
  },
});

export default Autocomplete;
