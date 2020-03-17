using System;

namespace ButtonBoxVjoyFeeder
{
    internal static class Errors
    {
        public static Exception Busy(uint id)
        {
            return new Exception($"vJoy Device {id} is already owned by another feeder");
        }

        public static Exception NotFound(uint id)
        {
            return new Exception($"vJoy Device {id} is not installed");
        }

        public static Exception UnknownError(uint id)
        {
            return new Exception($"vJoy Device {id} is not invalid");
        }

        internal static Exception CouldNotAcquite(uint id)
        {
            return new Exception($"Could not acquire device {id}");
        }

        public static Exception DriverMismatch(uint dllVersion, uint driverVersion)
        {
            return new Exception($"vJoy library version ({dllVersion} is not compatible with vJoy driver version {driverVersion}");
        }
    }
}