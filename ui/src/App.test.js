import { render, screen } from '@testing-library/react';
import App from './App';

test('renders src/App.js', () => {
  render(<App />);
  const element = screen.getByText(/src\/App.js/i);
  expect(element).toBeInTheDocument();
});
