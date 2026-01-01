import React from 'react';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler,
} from 'chart.js';
import { Line } from 'react-chartjs-2';
import { TrendPoint } from '../../store/slices/insightsSlice';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
);

interface TrendChartProps {
  data: TrendPoint[];
}

const TrendChart: React.FC<TrendChartProps> = ({ data }) => {
  if (!data || data.length === 0) {
    return (
      <div className="h-full flex items-center justify-center text-midnight-400">
        No trend data available
      </div>
    );
  }

  const chartData = {
    labels: data.map(point => {
      const date = new Date(point.date);
      return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
    }),
    datasets: [
      {
        label: 'Net Worth',
        data: data.map(point => point.value),
        fill: true,
        borderColor: 'rgba(45, 79, 255, 1)',
        backgroundColor: (context: any) => {
          const chart = context.chart;
          const { ctx, chartArea } = chart;
          if (!chartArea) return 'rgba(45, 79, 255, 0.1)';

          const gradient = ctx.createLinearGradient(0, chartArea.top, 0, chartArea.bottom);
          gradient.addColorStop(0, 'rgba(45, 79, 255, 0.3)');
          gradient.addColorStop(1, 'rgba(45, 79, 255, 0)');
          return gradient;
        },
        borderWidth: 2,
        pointRadius: 0,
        pointHoverRadius: 6,
        pointHoverBackgroundColor: 'rgba(45, 79, 255, 1)',
        pointHoverBorderColor: '#ffffff',
        pointHoverBorderWidth: 2,
        tension: 0.4,
      },
    ],
  };

  const options = {
    responsive: true,
    maintainAspectRatio: false,
    interaction: {
      intersect: false,
      mode: 'index' as const,
    },
    plugins: {
      legend: {
        display: false,
      },
      tooltip: {
        backgroundColor: 'rgba(12, 13, 31, 0.9)',
        titleColor: '#ffffff',
        bodyColor: 'rgba(142, 175, 255, 0.8)',
        borderColor: 'rgba(45, 79, 255, 0.3)',
        borderWidth: 1,
        padding: 12,
        cornerRadius: 8,
        displayColors: false,
        callbacks: {
          label: (context: any) => {
            return `$${context.parsed.y.toLocaleString()}`;
          },
        },
      },
    },
    scales: {
      x: {
        grid: {
          display: false,
        },
        ticks: {
          color: 'rgba(142, 175, 255, 0.5)',
          font: {
            family: 'Outfit',
            size: 11,
          },
          maxTicksLimit: 6,
        },
        border: {
          display: false,
        },
      },
      y: {
        grid: {
          color: 'rgba(45, 79, 255, 0.1)',
        },
        ticks: {
          color: 'rgba(142, 175, 255, 0.5)',
          font: {
            family: 'Outfit',
            size: 11,
          },
          callback: (value: any) => `$${(value / 1000).toFixed(0)}k`,
          maxTicksLimit: 5,
        },
        border: {
          display: false,
        },
      },
    },
  };

  return <Line data={chartData} options={options} />;
};

export default TrendChart;

