import './App.css';
import TickTable from '../components/tick_table';

export default function App() {
  return (
    <div className="App">
      <h1>Exchange price feed dashboard</h1>
      <TickTable />
    </div>
  )
};