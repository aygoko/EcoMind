import React from 'react';
import { View, Text, StyleSheet, Pressable } from 'react-native';
import Svg, { Path } from 'react-native-svg';

const EcoCard = ({ title, description, icon, onPress }) => {
  const TaskIcon = () => (
    <Svg width={40} height={40} viewBox="0 0 24 24" fill="none">
      <Path d="M19 21l-9-9-9 9" stroke="#2ecc71" strokeWidth={2} />
    </Svg>
  );

  return (
    <Pressable onPress={onPress} style={styles.card}>
      {icon === 'task' && <TaskIcon />}
      <Text style={styles.title}>{title}</Text>
      <Text style={styles.description}>{description}</Text>
    </Pressable>
  );
};

const styles = StyleSheet.create({
  card: {
    backgroundColor: 'white',
    borderRadius: 10,
    padding: 15,
    elevation: 3,
    width: '48%',
  },
  title: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 5,
    color: '#2c3e50',
  },
  description: {
    color: '#666',
  },
});

export default EcoCard;