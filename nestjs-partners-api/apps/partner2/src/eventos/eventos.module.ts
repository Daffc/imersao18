import { Module } from '@nestjs/common';
import { EventsCoreModule } from '@app/core/events/events-core.module';
import { EventosController } from './eventos.controller';

@Module({
  controllers: [EventosController],
  imports: [EventsCoreModule],
})
export class EventosModule { }
