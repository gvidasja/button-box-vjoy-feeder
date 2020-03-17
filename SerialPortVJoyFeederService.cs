using System;
using System.IO.Ports;
using System.Threading.Tasks;

namespace ButtonBoxVjoyFeeder
{
    public class SerialPortVJoyFeederService
    {
        private bool _running = true;
        private SerialPort _serialPort;
        private ButtonControllerPressSupervisor _buttonController;

        public void Start()
        {
            _serialPort = new SerialPort()
            {
                PortName = "COM15",
                BaudRate = 9600,
                DtrEnable = true
            };

            _serialPort.Open();

            _buttonController = new ButtonControllerPressSupervisor(VJoyAdapter.Init(), TimeSpan.FromMilliseconds(10));

            Task.Run(() =>
            {
                while (_running)
                {
                    var serialString = _serialPort.ReadLine();
                    var action = byte.Parse(serialString.Substring(0, 1)) > 0;
                    var button = byte.Parse(serialString.Substring(1));

                    _buttonController.SetButton(button, action);
                }
            });
        }

        public void Stop()
        {
            _running = false;
            _serialPort.Close();
            _buttonController.Dispose();
        }
    }
}