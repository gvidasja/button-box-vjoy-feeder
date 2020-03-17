using System;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace ButtonBoxVjoyFeeder
{
    public class ButtonControllerPressSupervisor : IButtonController
    {
        private readonly IDictionary<uint, DateTime> _pressMap;
        private readonly IButtonController _controller;
        private readonly TimeSpan _debounce;

        public ButtonControllerPressSupervisor(IButtonController controller, TimeSpan debounce)
        {
            _pressMap = new Dictionary<uint, DateTime>();
            _controller = controller;
            _debounce = debounce;
        }

        public void Dispose()
        {
            _controller.Dispose();
        }

        public Task SetButton(uint id, bool state)
        {
            var now = DateTime.UtcNow;

            if (state)
            {
                _pressMap[id] = now;
                return _controller.SetButton(id, state);
            }
            else if (_pressMap.TryGetValue(id, out var value) && now.Subtract(value) < _debounce)
            {
                return Delay(_debounce, () => _controller.SetButton(id, state));                
            }
            else
            {
                return _controller.SetButton(id, state);
            }
        }

        private async Task Delay(TimeSpan delay, Action action)
        {
            await Task.Delay(delay).ConfigureAwait(false);
            action();
        }
    }
}