using System;
using System.Threading.Tasks;
using vJoyInterfaceWrap;

namespace ButtonBoxVjoyFeeder
{
    public class VJoyAdapter : IDisposable, IButtonController
    {
        private readonly vJoy _joystick;
        private readonly uint _deviceId;

        private VJoyAdapter(uint id)
        {
            _deviceId = id;
            _joystick = new vJoy();
        }

        public static VJoyAdapter Init(uint id = 1)
        {
            return new VJoyAdapter(id).InitJoystick();
        }

        public void Dispose()
        {
            _joystick.RelinquishVJD(_deviceId);
        }

        public Task SetButton(uint buttonId, bool state)
        {
            _joystick.SetBtn(state, _deviceId, buttonId);

            return Task.CompletedTask;
        }

        private VJoyAdapter InitJoystick()
        {

            ValidateJoystick(_deviceId);

            if (!_joystick.AcquireVJD(_deviceId))
            {
                throw Errors.CouldNotAcquite(_deviceId);
            }
            else
            {
                Console.WriteLine($"Acquired: vJoy device {_deviceId}");
            }

            return this;
        }

        private void ValidateJoystick(uint id)
        {
            if (!_joystick.vJoyEnabled())
            {
                throw new Exception("vJoy is not enabled");
            }

            VjdStat status = _joystick.GetVJDStatus(id);

            switch (status)
            {
                case VjdStat.VJD_STAT_OWN:
                case VjdStat.VJD_STAT_FREE:
                    break;
                case VjdStat.VJD_STAT_BUSY:
                    throw Errors.Busy(id);
                case VjdStat.VJD_STAT_MISS:
                    throw Errors.NotFound(id);
                case VjdStat.VJD_STAT_UNKN:
                default:
                    throw Errors.UnknownError(id);
            };
        }
    }
}