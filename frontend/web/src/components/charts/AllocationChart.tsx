import React from 'react';
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js';
import { Doughnut } from 'react-chartjs-2';
import { AllocationItem } from '../../store/slices/insightsSlice';

ChartJS.register(ArcElement, Tooltip, Legend);

interface AllocationChartProps {
  data: AllocationItem[];
}

const defaultColors = [
  'rgba(45, 79, 255, 0.8)',   // midnight-500
  'rgba(16, 185, 129, 0.8)', // emerald
  'rgba(139, 92, 246, 0.8)', // violet
  'rgba(245, 158, 11, 0.8)', // amber
  'rgba(244, 63, 94, 0.8)',  // coral
  'rgba(142, 175, 255, 0.8)', // midnight-300
];

const AllocationChart: React.FC<AllocationChartProps> = ({ data }) => {
  if (!data || data.length === 0) {
    return (
      <div className="h-full flex items-center justify-center text-midnight-400">
        No allocation data available
      </div>
    );
  }

  const chartData = {
    labels: data.map(item => item.category),
    datasets: [
      {
        data: data.map(item => item.percentage),
        backgroundColor: data.map((item, index) => item.color || defaultColors[index % defaultColors.length]),
        borderColor: 'rgba(12, 13, 31, 1)',
        borderWidth: 3,
        hoverOffset: 8,
      },
    ],
  };

  const options = {
    responsive: true,
    maintainAspectRatio: false,
    cutout: '65%',
    plugins: {
      legend: {
        position: 'right' as const,
        labels: {
          color: 'rgba(142, 175, 255, 0.8)',
          padding: 16,
          font: {
            family: 'Outfit',
            size: 12,
          },
          usePointStyle: true,
          pointStyle: 'circle',
        },
      },
      tooltip: {
        backgroundColor: 'rgba(12, 13, 31, 0.9)',
        titleColor: '#ffffff',
        bodyColor: 'rgba(142, 175, 255, 0.8)',
        borderColor: 'rgba(45, 79, 255, 0.3)',
        borderWidth: 1,
        padding: 12,
        cornerRadius: 8,
        displayColors: true,
        callbacks: {
          label: (context: any) => {
            const value = context.parsed;
            const amount = data[context.dataIndex]?.amount || 0;
            return ` ${value.toFixed(1)}% ($${amount.toLocaleString()})`;
          },
        },
      },
    },
  };

  return <Doughnut data={chartData} options={options} />;
};

export default AllocationChart;

