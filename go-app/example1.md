Анализ проекта "OverwatchArcade.API"

Общее количество ошибок: 5

Архитектурных нарушений: 2

Несоответствий стандартам: 3

### Архитектурное нарушение 1
> "AutoMapperProfile.cs" (строка 27, символ 13)
> Необходимо использовать интерфейсы вместо конкретных реализаций коллекций.
```csharp
CreateMap<Daily, DailyDto>()
    .ForMember(dest => dest.Modes, opt => opt.MapFrom(src => src.TileModes));
```
> Предложенное исправление:
```csharp
CreateMap<IDaily, DailyDto>()
    .ForMember(dest => dest.Modes, opt => opt.MapFrom(src => src.TileModes));
```

### Архитектурное нарушение 2
> "AutoMapperProfile.cs" (строка 33, символ 13)
> Необходимо использовать интерфейсы вместо конкретных реализаций коллекций.
```csharp
CreateMap<TileMode, TileModeDto>()
    .ForMember(dest => dest.Name, opt => opt.MapFrom(src => src.ArcadeMode.Name))
    // ... другие члены ...
```
> Предложенное исправление:
```csharp
CreateMap<ITileMode, TileModeDto>()
    .ForMember(dest => dest.Name, opt => opt.MapFrom(src => src.ArcadeMode.Name))
    // ... другие члены ...
```

### Несоответствие стандартам 1
> "AutoMapperProfile.cs" (строка 22, символ 5)
> Необходимо проверить аргументы внутри методов.
```csharp
CreateMap<Daily, DailyDto>();
```
> Предложенное исправление:
```csharp
if (typeof(Daily) != null && typeof(DailyDto) != null)
{
    CreateMap<Daily, DailyDto>();
}
```

### Несоответствие стандартам 2
> "AutoMapperProfile.cs" (строка 24, символ 5)
> Необходимо проверить аргументы внутри методов.
```csharp
CreateMap<CreateDailyDto, Daily>();
```
> Предложенное исправление:
```csharp
if (typeof(CreateDailyDto) != null && typeof(Daily) != null)
{
    CreateMap<CreateDailyDto, Daily>();
}
```

### Несоответствие стандартам 3
> "AutoMapperProfile.cs" (строка 38, символ 5)
> Необходимо проверить аргументы внутри методов.
```csharp
CreateMap<ContributorProfileDto, ContributorProfile>();
```
> Предложенное исправление:
```csharp
if (typeof(ContributorProfileDto) != null && typeof(ContributorProfile) != null)
{
    CreateMap<ContributorProfileDto, ContributorProfile>();
}
```

### Рекомендации по NuGet-зависимостям
Все пакеты должны быть обновлены до актуальных версий. Уязвимые транзитивные зависимости должны быть обновлены или включены в проект. Отсутствие лишних зависимостей и абсолютных путей.

### Рекомендации по коду
Удалите незавершённые TODO, закомментированный/неиспользуемый код, устаревший [Obsolete] код. Комментарии должны быть у ключевых элементов. Исключите: неиспользуемые переменные, дублирование логов исключений, некорректное объединение строк, избыточные null-проверки.

### Рекомендации по LINQ
Применяйте Chunk() вместо Skip().Take(). Для пользовательских типов с Union(), Distinct() и т.п. переопределите Equals() и GetHashCode(). Избегайте лишних вызовов ToList() и Distinct() после Union().

### Рекомендации по Entity Framework
Используйте асинхронные методы материализации данных. Не удаляйте сущности в циклах, группируйте SaveChangesAsync(). Фильтрацию данных выполняйте на стороне БД. Применяйте AddAsync()/AddRangeAsync() только с SequenceHiLo.
