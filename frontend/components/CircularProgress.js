import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { LineChart } from 'react-native-chart-kit';

const CircularProgress = ({ progress, goal }) => {
  const data = {
    labels: ['Прогресс', 'Осталось'],
    datasets: [
      {
        data: [progress, goal - progress],
        color: (opacity = 1) => `rgba(46, 204, 113, ${opacity})`,
      },
    ],
  };

  return (
    <View style={styles.chartContainer}>
      <LineChart
        data={data}
        width={250}
        height={220}
        chartConfig={{
          backgroundColor: '#e2e2e2',
          backgroundGradientFrom: '#ffffff',
          backgroundGradientTo: '#ffffff',
          decimalPlaces: 0,
          color: (opacity = 1) => `rgba(0, 0, 0, ${opacity})`,
        }}
        bezier
        style={{
          marginVertical: 8,
          borderRadius: 16,
        }}
      />
      <Text style={styles.progressText}>
        {progress}/{goal}
      </Text>
    </View>
  );
};

const styles = StyleSheet.create({
  chartContainer: {
    alignItems: 'center',
    marginBottom: 20,
  },
  progressText: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#2ecc71',
  },
});

export default CircularProgress;