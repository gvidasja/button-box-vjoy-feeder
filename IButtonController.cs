using System;
using System.Threading.Tasks;

namespace ButtonBoxVjoyFeeder
{
    public interface IButtonController : IDisposable
    {
        Task SetButton(uint id, bool state);
    }
}